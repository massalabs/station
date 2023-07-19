import { FiDownload } from 'react-icons/fi';
import { Plugin } from '@massalabs/react-ui-kit';
export function PluginCardLoading() {
  const argsStore = {
    preIcon: <img src="" alt="plugin-logo" />,
    topAction: <FiDownload />,
    title: '',
    subtitle: '',
    subtitleIcon: null,
    content: '',
  };

  return (
    <>
      <div className="w-fit h-fit blur-sm animate-pulse">
        <Plugin {...argsStore} />
      </div>
    </>
  );
}
