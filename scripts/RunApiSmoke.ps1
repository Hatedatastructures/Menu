param(
    [Parameter(Mandatory = $true)]
    [string]$Executable,
    [Parameter(Mandatory = $true)]
    [string]$Config,
    [int]$Port = 8080
)

$ExecutablePath = (Resolve-Path -LiteralPath $Executable).Path
$ConfigPath = (Resolve-Path -LiteralPath $Config).Path
$WorkingDirectory = (Get-Location).Path
$Process = Start-Process -FilePath $ExecutablePath `
    -ArgumentList @('--config', $ConfigPath) `
    -WorkingDirectory $WorkingDirectory `
    -PassThru

try {
    $BaseUrl = "http://127.0.0.1:$Port"
    $Ready = $false
    for ($Attempt = 0; $Attempt -lt 50; $Attempt++) {
        if ($Process.HasExited) {
            throw "MenuServer exited before health check with code $($Process.ExitCode)"
        }
        try {
            $HealthResponse = Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/healthz" -TimeoutSec 1
            if ($HealthResponse.StatusCode -eq 200) {
                $Ready = $true
                break
            }
        } catch {
        }
        Start-Sleep -Milliseconds 100
    }
    if (-not $Ready) {
        throw 'MenuServer health check did not become ready'
    }

    $HealthResponse = Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/healthz" -TimeoutSec 5
    $ReadyResponse = Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/readyz" -TimeoutSec 5
    $RecipesResponse = Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/api/v1/recipes" -TimeoutSec 5
    $IngredientsResponse = Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/api/v1/ingredients" -TimeoutSec 5
    $Recipes = $RecipesResponse.Content | ConvertFrom-Json
    if ($Recipes.Count -eq 0) {
        throw 'Recipe response did not contain any recipe'
    }
    $RecipeId = [uri]::EscapeDataString($Recipes[0].id)
    $RecipeResponse = Invoke-WebRequest -UseBasicParsing -Uri "$BaseUrl/api/v1/recipes/$RecipeId" -TimeoutSec 5
    $RecommendationBody = '{"servings":2,"availableMinutes":35,"cuisines":["中餐"],"pantryIngredientIds":["ingredient.rice"],"cookware":["炒锅"]}'
    $RecommendationResponse = Invoke-WebRequest -UseBasicParsing `
        -Method Post -ContentType 'application/json' -Body $RecommendationBody `
        -Uri "$BaseUrl/api/v1/recommendations/tonight" -TimeoutSec 5
    if ($HealthResponse.StatusCode -ne 200 -or
        $ReadyResponse.StatusCode -ne 200 -or
        $RecipesResponse.StatusCode -ne 200 -or
        $IngredientsResponse.StatusCode -ne 200 -or
        $RecipeResponse.StatusCode -ne 200 -or
        $RecommendationResponse.StatusCode -ne 200) {
        throw "Unexpected smoke status: health=$($HealthResponse.StatusCode), ready=$($ReadyResponse.StatusCode), recipes=$($RecipesResponse.StatusCode), ingredients=$($IngredientsResponse.StatusCode), recipe=$($RecipeResponse.StatusCode), recommendations=$($RecommendationResponse.StatusCode)"
    }
    $Ingredients = $IngredientsResponse.Content | ConvertFrom-Json
    if ($Ingredients.Count -eq 0 -or [string]::IsNullOrWhiteSpace($Ingredients[0].name)) {
        throw 'Ingredient response did not contain display fields'
    }

    [PSCustomObject]@{
        ProcessId = $Process.Id
        HealthStatus = $HealthResponse.StatusCode
        ReadyStatus = $ReadyResponse.StatusCode
        RecipesStatus = $RecipesResponse.StatusCode
        IngredientsStatus = $IngredientsResponse.StatusCode
        RecipeStatus = $RecipeResponse.StatusCode
        RecommendationStatus = $RecommendationResponse.StatusCode
        RecipesBytes = $RecipesResponse.RawContentLength
        IngredientsBytes = $IngredientsResponse.RawContentLength
    } | ConvertTo-Json -Compress
} finally {
    if ($null -ne $Process -and -not $Process.HasExited) {
        Stop-Process -Id $Process.Id -Force
        $Process.WaitForExit()
    }
}
