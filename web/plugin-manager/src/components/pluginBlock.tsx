import React, { useEffect, useMemo, useRef } from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
export type PluginProps = {
    name: string;
    logo: string;
    description: string;
    version: string;
    online: boolean;
    updateAvailable: boolean;
    id: number;
};

function PluginBlock(props: PluginProps) {
    //UseMemo props
    // useMemo(() => {
    //   console.log("props UseMemo", props);
    // }, [props]);
    // Each Rerender we update fetch plugins data
    useEffect(() => {
        fetchPluginInfo();
        console.log("props UseEffect", props);
    }, [props]);
    //UseRef props
    const propsRef = useRef(props);

    // fetch info from plugin
    function fetchPluginInfo() {
        fetch(`${window.location.hostname}/thyra/plugin-manager/${props.id}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Success:", data);
                if (data !== undefined) {
                    propsRef.current = data;
                }
            });
    }

    function launchAndStopPlugins() {
        // fetch info from plugin
        fetch(`${window.location.hostname}/thyra/plugin-manager/${props.id}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Success:", data);
            });
        console.log(props);
        // Launch plugin
        return launchPlugins();
    }

    // Launch plugin
    function launchPlugins() {
        fetch(`${window.location.hostname}/thyra/plugin-manager/${props.id}/execute`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                command: "start",
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Success:", data);
                if (data !== undefined) {
                    propsRef.current.online = true;
                }
            });
    }
    // Stop plugin
    function stopPlugins() {
        fetch(`${window.location.hostname}/thyra/plugin-manager/${props.id}/execute`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                command: "stop",
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Success:", data);
                if (data !== undefined) {
                    propsRef.current.online = false;
                }
            });
    }
    // restart plugin
    function restartPlugins() {
        fetch(`${window.location.hostname}/thyra/plugin-manager/${props.id}/execute`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                command: "restart",
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Success:", data);
                if (data !== undefined) {
                    propsRef.current.online = true;
                }
            });
    }

    function minimizeString(str: string, length: number) {
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    function playStatus() {
        return props.online ? "w-6 h-6 text-green-500" : "w-6 h-6 text-red-500";
    }

    function updateStatus() {
        return props.updateAvailable ? "w-6 h-6 text-yellow-500" : "w-6 h-6 text-green-500";
    }

    function toggleStatus() {
        return props.online ? "toggle toggle-success" : "toggle toggle-success checked";
    }

    return (
        <section className="bg-slate-800 h-48 max-w-lg w-96 p-3 m-4 rounded-2xl">
            <div className=" flex-row h-full text-white ">
                <div className="flex">
                    <img className="w-10 h-10 pt-3 mx-2" src={props.logo} alt="Plugin Logo" />
                    <div className="w-full">
                        <h1 className="font-bold">{minimizeString(props.name, 90)}</h1>
                        <p className="font-light max-sm:text-sm">
                            {minimizeString(props.description, 100)}
                        </p>
                    </div>
                </div>
                <div className="content-center flex-wrap flex ">
                    <div className="flex w-full pt-7 justify-around items-center">
                        <p className="font-light">V: {props.version}</p>

                        <input
                            type="checkbox"
                            className="toggle toggle-success"
                            checked={props.online}
                            onChange={launchAndStopPlugins}
                        />

                        <button>
                            <PlayCircleIcon className={playStatus()} />
                        </button>
                        <button>
                            <ArrowPathIcon className={updateStatus()} />
                        </button>
                        <button>
                            <TrashIcon className="w-6 h-6 " />
                        </button>
                    </div>
                </div>
            </div>
        </section>
    );
}

export default PluginBlock;
