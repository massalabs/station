import React from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
export type PluginProps = {
    name: string;
    logo: string;
    description: string;
    version: string;
    online: boolean;
    updateAvailable: boolean;
};

function PluginBlock(props: PluginProps) {
    function minimizeString(str: string, length: number) {
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    function playStatus() {
        if (props.online) {
            return <PlayCircleIcon className="w-6 h-6 text-green-500" />;
        } else {
            return <PlayCircleIcon className="w-6 h-6 text-red-500" />;
        }
    }
    function updateStatus() {
        if (props.updateAvailable) {
            return (
                <div>
                    <ArrowPathIcon className="w-6 h-6 text-yellow-500" />
                </div>
            );
        } else {
            return (
                <div>
                    <ArrowPathIcon className="w-6 h-6 text-green-500" />
                </div>
            );
        }
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

                        <input type="checkbox" className="toggle toggle-success " />

                        <button>{updateStatus()}</button>
                        <button>{playStatus()}</button>
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
