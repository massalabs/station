import Toggle from 'react-toggle';
import moon from '../assets/pictos/moon.svg';
import sun from '../assets/pictos/sunWhite.svg';
import { UIStore } from '../store/UIStore';
import 'react-toggle/style.css';
import '../index.css';

const toggleTheme = () => {
  UIStore.useState((s) => s.theme);
  const handleChange = () => {
    UIStore.update((s) => {
      s.theme = s.theme == 'dark' ? 'light' : 'dark';
    });
  };
  return (
    <Toggle
      defaultChecked={UIStore.useState((s) =>
        s.theme == 'light' ? true : false,
      )}
      icons={{
        checked: <img src={sun} />,
        unchecked: <img src={moon} />,
      }}
      className="custom-classname"
      onChange={handleChange}
    />
  );
};

export default toggleTheme;
