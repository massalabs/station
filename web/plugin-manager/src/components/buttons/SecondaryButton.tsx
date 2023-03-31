import React from "react";
import { UIStore } from "../../store/UIStore";

type Props = {
    label: string;
    onClick: () => void;
    width?: string;
    isDisabled?: boolean;
    type?: "button" | "submit" | "reset";
};

const SecondaryButton = (props: Props) => {
    return (
        <button
            type={props.type ?? "button"}
            onClick={props.onClick}
            className={
                "flex flex-row justify-center items-center gap-2 h-12 border-[1px] border-solid border-border rounded-md " +
                (props.width ?? " w-28") +
                (props.isDisabled
                    ? " bg-disabledButton cursor-not-allowed"
                    : " bg-secondaryButton hover:bg-hoverPrimaryButton cursor-pointer")
                }
        >
            <p className={"text-font"}>{props.label}</p>
        </button>
    );
};

export default SecondaryButton;
