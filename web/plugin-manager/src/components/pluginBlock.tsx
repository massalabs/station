import React from "react";
import { ArrowPathIcon, TrashIcon, PlayCircleIcon } from "@heroicons/react/24/outline";
export type PluginProps = {
    name: string;
    logo: string;
    description: string;
    version: string;
    online: boolean;
};

function PluginBlock(props: PluginProps) {
    function minimizeString(str: string, length: number) {
        if (str.length > length) {
            return str.substring(0, length) + "...";
        } else {
            return str;
        }
    }

    return (
        <div className="bg-slate-800 max-w-lg w-96 p-3 m-4 text-white rounded-md">
            <div className="flex flex-row ">
                <img className="w-10 h-10 mx-2 top" src={props.logo} alt="Plugin Logo" />
                <div className="flex flex-col">
                  <h1 className="font-bold text-center">{minimizeString(props.name, 90)}</h1>
                  <p className="font-light">{minimizeString(props.description, 100)}</p>
                </div>
            </div>
            <div className="inline-flex h-full w-full self-end pt-7">
                <div className="">
                    <div className="inline-flex mx-auto ">
                        <p className="font-light ml-3">V: {props.version}</p>
                        <input type="checkbox" className="toggle toggle-success ml-3" />
                        <button>
                            <PlayCircleIcon className="w-6 h-6 ml-3" />
                        </button>
                        <button>
                            <ArrowPathIcon className="w-6 h-6 ml-3" />
                        </button>
                        <button>
                            <TrashIcon className="w-6 h-6 ml-3" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default PluginBlock;
