const BaseUrl = process.env.MENU_BENCHMARK_URL || "http://127.0.0.1:8080";
const TotalRequests = Number(process.env.MENU_BENCHMARK_REQUESTS || 512);
const Concurrency = Number(process.env.MENU_BENCHMARK_CONCURRENCY || 16);

if (!Number.isInteger(TotalRequests) || TotalRequests < 1 ||
    !Number.isInteger(Concurrency) || Concurrency < 1) {
    throw new Error("MENU_BENCHMARK_REQUESTS and MENU_BENCHMARK_CONCURRENCY must be positive integers");
}

const Durations = [];
let ErrorCount = 0;
let NextRequest = 0;
const StartedAt = performance.now();

async function RunWorker() {
    while (true) {
        const RequestNumber = NextRequest++;
        if (RequestNumber >= TotalRequests) {
            return;
        }
        const RequestStartedAt = performance.now();
        try {
            const Response = await fetch(`${BaseUrl}/api/v1/recipes?limit=20`, {
                headers: {Accept: "application/json"},
            });
            await Response.arrayBuffer();
            if (!Response.ok) {
                ErrorCount += 1;
            }
        } catch {
            ErrorCount += 1;
        }
        Durations.push(performance.now() - RequestStartedAt);
    }
}

await Promise.all(Array.from(
    {length: Math.min(Concurrency, TotalRequests)},
    () => RunWorker(),
));

Durations.sort((Left, Right) => Left - Right);
const Percentile = (Ratio) => {
    const Index = Math.min(Durations.length - 1,
        Math.ceil(Ratio * Durations.length) - 1);
    return Durations[Index];
};
const ElapsedMilliseconds = performance.now() - StartedAt;

console.log(JSON.stringify({
    Url: BaseUrl,
    Requests: TotalRequests,
    Concurrency,
    Errors: ErrorCount,
    ErrorRate: ErrorCount / TotalRequests,
    ElapsedMilliseconds: Number(ElapsedMilliseconds.toFixed(2)),
    RequestsPerSecond: Number((TotalRequests / (ElapsedMilliseconds / 1000)).toFixed(2)),
    P50Milliseconds: Number(Percentile(0.5).toFixed(2)),
    P95Milliseconds: Number(Percentile(0.95).toFixed(2)),
    P99Milliseconds: Number(Percentile(0.99).toFixed(2)),
}));

if (ErrorCount !== 0) {
    process.exitCode = 1;
}
