import { useState } from "react";

import PrimaryButton from "./buttons/PrimaryButton";
import { BarLoader } from "react-spinners";

type Props = {
    callbackToParent: (data: string) => void;
    label?: string;
    placeholder: string;
    buttonValue: string;
    error?: string;
};

const LabelButton = (props: Props) => {
    const [value, setValue] = useState("");
    const [isPending, setIsPending] = useState(false);

    function handleInputValueChange(event: any) {
        setValue(event.target.value);
    }

    function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();
        props.callbackToParent(value);
        setIsPending(true);
        setTimeout(() => {
            setIsPending(false);
        }, 4000);
    }

    return (
        <form onSubmit={handleSubmit} className="">
            {props.error && <p className="text-red-500">{props.error}</p>}
            <input
                type="text"
                className="text-font w-full mb-4 rounded-md bg-primaryBG"
                placeholder={props.placeholder}
                onChange={handleInputValueChange}
            />
            {!isPending ? (
                <PrimaryButton
                    label={props.buttonValue}
                    type="submit"
                    width={" w-full"}
                    onClick={() => props.callbackToParent(value)}
                />
            ) : (
                <BarLoader width={"100%"} color="hsl(var(--twc-brand))" />
            )}
        </form>
    );
};

export default LabelButton;
