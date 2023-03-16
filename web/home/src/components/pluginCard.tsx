import React, { MouseEventHandler, MouseEvent} from 'react';
import Arrow6 from '../assets/pictos/arrow6.svg';
import { PluginHomePage } from '../../../shared/interfaces/IPlugin';
type Props = {
  plugin: PluginHomePage;
  handleOpenPlugin: (event: MouseEvent<HTMLDivElement>) => void;
  key: string;
};

export const PluginCard = (props: Props) => {
  const handleCardClick = () => {
    props.handleOpenPlugin(props.plugin.name);
  };
  return (
    <div
      onClick={handleCardClick}
      className="flex flex-col justify-center items-start ml-2 p-5 gap-4 w-64 h-72 
                    border-[1px] border-solid border-black rounded-2xl bg-bgCard cursor-pointer"
    >
      <img
        src={props.plugin.logo}
        alt="Album"
        className="rounded-3xl w-10 h-10"
      />
      <div className="flex flex-col gap-2">
        <h2 className="label2 text-font">{props.plugin.name}</h2>
        <p className="text2 text-font overflow-hidden whitespace-pre-wrap max-w-full">
          {props.plugin.description}
        </p>
      </div>
      <img src={Arrow6} alt="Album" className="w-6 h-6" />
    </div>
  );
};
