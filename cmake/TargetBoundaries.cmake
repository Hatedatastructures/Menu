function(CollectMenuTargetClosure Target OutVariable)
    set(Pending "${Target}")
    set(Visited)
    set(Closure)

    while(Pending)
        list(POP_FRONT Pending Current)
        if(NOT TARGET "${Current}")
            continue()
        endif()

        list(FIND Visited "${Current}" AlreadyVisited)
        if(NOT AlreadyVisited EQUAL -1)
            continue()
        endif()
        list(APPEND Visited "${Current}")

        foreach(Property LINK_LIBRARIES INTERFACE_LINK_LIBRARIES)
            get_target_property(Dependencies "${Current}" "${Property}")
            if(NOT Dependencies OR Dependencies MATCHES "-NOTFOUND$")
                continue()
            endif()
            foreach(Dependency IN LISTS Dependencies)
                if(Dependency MATCHES "^\$<")
                    continue()
                endif()
                if(TARGET "${Dependency}")
                    list(APPEND Closure "${Dependency}")
                    list(APPEND Pending "${Dependency}")
                endif()
            endforeach()
        endforeach()
    endwhile()

    list(REMOVE_DUPLICATES Closure)
    set(${OutVariable} "${Closure}" PARENT_SCOPE)
endfunction()

function(AssertMenuTargetBoundaries)
    set(RequiredTargets
        MenuFoundation
        MenuDomain
        MenuApplication
        MenuInfrastructure
        MenuRuntime
        MenuTransportCore
        MenuTransport
        MenuApi
        MenuComposition
        MenuServer
    )

    foreach(Target IN LISTS RequiredTargets)
        if(NOT TARGET "${Target}")
            message(FATAL_ERROR "Required Menu target is missing: ${Target}")
        endif()
    endforeach()

    foreach(Target IN LISTS RequiredTargets)
        get_target_property(DirectDependencies "${Target}" LINK_LIBRARIES)
        if(DirectDependencies AND NOT DirectDependencies MATCHES "-NOTFOUND$")
            list(FIND DirectDependencies "${Target}" SelfDependency)
            if(NOT SelfDependency EQUAL -1)
                message(FATAL_ERROR "Menu target directly depends on itself: ${Target}")
            endif()
        endif()
    endforeach()

    CollectMenuTargetClosure(MenuTransportCore TransportCoreClosure)
    list(FIND TransportCoreClosure MenuApi TransportCoreApi)
    if(NOT TransportCoreApi EQUAL -1)
        message(FATAL_ERROR
            "MenuTransportCore closure must not contain MenuApi: ${TransportCoreClosure}")
    endif()

    CollectMenuTargetClosure(MenuTransport TransportClosure)
    list(FIND TransportClosure MenuApi TransportApi)
    if(NOT TransportApi EQUAL -1)
        message(FATAL_ERROR
            "MenuTransport closure must not contain MenuApi: ${TransportClosure}")
    endif()

    CollectMenuTargetClosure(MenuApi ApiClosure)
    list(FIND ApiClosure MenuTransport ApiTransport)
    if(NOT ApiTransport EQUAL -1)
        message(FATAL_ERROR
            "MenuApi closure must not contain concrete MenuTransport: ${ApiClosure}")
    endif()
    list(FIND ApiClosure MenuTransportCore ApiTransportCore)
    if(ApiTransportCore EQUAL -1)
        message(FATAL_ERROR
            "MenuApi must consume MenuTransportCore: ${ApiClosure}")
    endif()

    CollectMenuTargetClosure(MenuRuntime RuntimeClosure)
    foreach(ForbiddenRuntimeDependency IN ITEMS
            MenuApi MenuApplication MenuInfrastructure MenuTransport)
        list(FIND RuntimeClosure "${ForbiddenRuntimeDependency}" RuntimeDependencyIndex)
        if(NOT RuntimeDependencyIndex EQUAL -1)
            message(FATAL_ERROR
                "MenuRuntime must not depend on ${ForbiddenRuntimeDependency}: ${RuntimeClosure}")
        endif()
    endforeach()

    file(GLOB_RECURSE RuntimeSources
        "${PROJECT_SOURCE_DIR}/server/Runtime/*.cpp"
        "${PROJECT_SOURCE_DIR}/server/Runtime/*.hpp")
    foreach(RuntimeSource IN LISTS RuntimeSources)
        file(READ "${RuntimeSource}" RuntimeSourceText)
        if(RuntimeSourceText MATCHES "#[ \\t]*include[ \\t]*[<\"](Api|Application|Infrastructure)/")
            message(FATAL_ERROR
                "MenuRuntime source must not include a business or storage module: ${RuntimeSource}")
        endif()
    endforeach()

    CollectMenuTargetClosure(MenuServer ServerClosure)
    list(FIND ServerClosure MenuComposition ServerComposition)
    if(ServerComposition EQUAL -1)
        message(FATAL_ERROR
            "MenuServer must compose MenuComposition: ${ServerClosure}")
    endif()

    message(STATUS "Menu target boundary check passed")
    message(STATUS "MenuTransportCore closure: ${TransportCoreClosure}")
    message(STATUS "MenuTransport closure: ${TransportClosure}")
    message(STATUS "MenuApi closure: ${ApiClosure}")
    message(STATUS "MenuRuntime closure: ${RuntimeClosure}")
    CollectMenuTargetClosure(MenuComposition CompositionClosure)
    message(STATUS "MenuComposition closure: ${CompositionClosure}")
    message(STATUS "MenuServer closure: ${ServerClosure}")
endfunction()
