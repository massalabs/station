import React, { useState } from "react";
import halfArrow from "../assets/pictos/halfArrow.svg";
import halfArrowWhite from "../assets/pictos/halfArrowWhite.svg";
import { UIStore } from "../store/UIStore";
import LabelButton from "./LabelButton";

type Props = {
    label2: string;
    text3: string;
    propsLabelButton: {
        callbackToParent: (data: string) => void;
        label: string;
        placeholder: string;
        buttonValue: string;
        error?: string;
    };
};

const SmallCardExtended = (props: Props) => {
    const isThemeLight = UIStore.useState((s) => (s.theme === "light"));
    const [isExtended, setIsExtended] = useState(false);
    const HandleOnClickExtend = () => {
        setIsExtended(!isExtended);
        //Change the className of the container to grow the height of the section
        return " h-56 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard";
    };
    const sectionClassName = () => {
        if (isExtended) {
            return "h-56 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard hover:bg-hoverbgCard";
        } else {
            return "h-28 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard hover:bg-hoverbgCard";
        }
    };
    const imgClassName = () => {
        if (isExtended) {
            return "animate-rotate180Smooth";
        } else {
            return "animate-rotateReverse180Smooth";
        }
    };
    return (
        <section className={sectionClassName()}>
            {/* Card Definition  */}
            <div className="flex flex-row " onClick={HandleOnClickExtend}>
                <div className="flex flex-col gap-3">
                    <p className="label2 text-font">{props.label2}</p>
                    <p className="text3 text-font">{props.text3}</p>
                </div>
                <div className="flex self-center mx-auto">
                    <img
                        className={"w-8 h-4 " + imgClassName()}
                        src={isThemeLight ? halfArrow : halfArrowWhite}
                        alt=""
                        onClick={HandleOnClickExtend}
                    />
                    {/* contains the icon to grow the container */}
                </div>
            </div>
            {isExtended && (
                <div className="flex flex-col mt-3">
                    <LabelButton {...props.propsLabelButton} />
                </div>
            )}
        </section>
    );
};

export default SmallCardExtended;
