import Toggle from "react-toggle";
import moon from "../assets/pictos/moon.svg";
import sun from "../assets/pictos/sunWhite.svg";
import { UIStore } from "../store/UIStore";
import "react-toggle/style.css";
import "../index.css";

type Props = {
    handleChange: () => void;
    checked: boolean;
};

const TogglePlugin = (props: Props) => {
    const handleChange = () => {
        props.handleChange();
    };
    return (
        <Toggle
            defaultChecked={props.checked}
            icons={false}
            className={`custom-TogglePlugin-${UIStore.useState((s) => s.theme)}`}
            onChange={handleChange}
        />
    );
};

export default TogglePlugin;
