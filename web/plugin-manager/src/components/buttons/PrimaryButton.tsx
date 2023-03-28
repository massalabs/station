import React from "react";
import { UIStore } from "../../store/UIStore";

type Props = {
    label: string;
    onClick: () => void;
    iconPathLight?: string;
    iconPathDark?: string;
    width?: string;
    type?: "button" | "submit" | "reset";
};

const PrimaryButton = (props: Props) => {
    const isThemeLight = UIStore.useState((s) => (s.theme == "light" ? true : false));
    return (
        <button
            type={props.type ?? "button"}
            onClick={props.onClick}
            className={
                "flex flex-row justify-center items-center gap-2 h-12  bg-primaryButton border-[1px] border-solid border-border rounded-md cursor-pointer" +
                (props.width ?? " w-28 ") +
                " hover:bg-hoverSecondaryButton "
            }
        >
            <p className="text-invertedfont">{props.label}</p>
            {props.iconPathLight ? (
                <img
                    className="w-2 h-2 mt-1"
                    src={isThemeLight ? props.iconPathLight : props.iconPathDark}
                    alt="No Icon Found"
                />
            ) : (
                <></>
            )}
        </button>
    );
};

export default PrimaryButton;
