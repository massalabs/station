import React from "react";
import Toggle from "react-toggle";
import moon from "../../public/pictos/moon.svg";
import sun from "../../public/pictos/sun.svg";
type Props = {
    state: boolean;
    setState: (state: boolean) => void;
};

const toggleTheme = (props: Props) => {
    const handleChange = () => {
        props.setState(!props.state);
    };
    return (
        <label>
            <Toggle
                defaultChecked={false}
                icons={{
                    checked: moon,
                    unchecked: sun,
                }}
                onChange={handleChange}
            />
            <span>Custom icons</span>
        </label>
    );
};

export default toggleTheme;
