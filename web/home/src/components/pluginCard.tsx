import Arrow6 from '../assets/pictos/arrow6.svg';
import ArrowWhite6 from '../assets/pictos/ArrowWhite6.svg';
import { PluginHomePage } from '../../../shared/interfaces/IPlugin';
import { UIStore } from '../store/UIStore';
type Props = {
  plugin: PluginHomePage;
  handleOpenPlugin: (pluginName: string) => void;
  key: string;
};

export const PluginCard = (props: Props) => {
  const handleCardClick = () => {
    props.handleOpenPlugin(props.plugin.name);
  };
  const Arrow = UIStore.useState((s) =>
    s.theme == 'light' ? Arrow6 : ArrowWhite6,
  );

  return (
    <div
      onClick={handleCardClick}
      className="flex flex-col justify-center items-start p-6 gap-4 w-72 h-56 
      border-[1px] border-solid border-border rounded-2xl cursor-pointer bg-bgCard hover:bg-hoverbgCard"
    >
      <img
        src={props.plugin.logo}
        alt="Album"
        className="rounded-3xl w-10 h-10"
      />
      <div className="flex flex-col gap-2 w-full">
        <h2 className={`label2 text-font h-8 minimize`}>{props.plugin.name}</h2>
        <p className={`text2 text-font max-w-full h-[68px] minimize`}>
          {props.plugin.description}
        </p>
      </div>
      <img src={Arrow} alt="Album" className="w-6 h-6" />
    </div>
  );
};
