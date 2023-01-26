import React from "react";
import PluginBlock from "../components/pluginBlock";
import { PluginProps } from "../components/pluginBlock";
import massaLogoLight from "../assets/MASSA_LIGHT_Detailed.png";
function Manager() {
    const mock: PluginProps = {
        name: "Plugin 1",
        logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description: "T Becarefull",
        version: "1.0.0",
        online: false,
        updateAvailable: true,
        id:1,
    };
    const mock2: PluginProps = {
        name: "Plugin 2",
        logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description: "This is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        online: false,
        updateAvailable: false,
        id: 2,
    };
    const mock3: PluginProps = {
        name: "Plugin 3",
        logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
        description: "This is a plugin Descriaaaaaaaaaaaaatio  aaaaaaaaaaa n  aaaaaaaaaaa BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
        version: "1.0.0",
        online: true,
        updateAvailable: true,
        id: 3,
    };
    const mock4: PluginProps = {
      name: "Plugin 4",
      logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
      description: "This is a plugin Deaaaaaaaaaacription BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
      version: "1.0.0",
      online: true,
      updateAvailable: false,
      id: 4,
  };
  const mock5: PluginProps = {
    name: "Plugin 5",
    logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
    description: "This is a plugin Description BessssssssssscarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
    version: "1.0.0",
    online: true,
    updateAvailable: true,
    id: 5,
};
const mock6: PluginProps = {
  name: "Plugin 5",
  logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/1/15/Audi_logo.svg/1280px-Audi_logo.svg.png",
  description: "This is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description BecarefullThis is a plugin Description Becarefull",
  version: "1.0.0",
  online: true,
  updateAvailable: true,
  id: 6,
};
    const mocks = [mock, mock2, mock3, mock4, mock5, mock6];
    const plugins = mocks.map((mock) => <PluginBlock {...mock} />);
    return (
        <>
            <div className="p-5 flex items-center">
                <img className="max-h-6" src={massaLogoLight} alt="Thyra Logo" />
                <h1 className="text-xl ml-6 font-bold text-white">Thyra</h1>
            </div>
            {/* FlexWrap is blocking align content in Plugin Block*/}
            <div className="flex flex-row flex-wrap mx-auto max-w-6xl justify-center content-center">
              {plugins}
            </div>
        </>
    );
}

export default Manager;
