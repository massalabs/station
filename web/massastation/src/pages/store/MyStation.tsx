import { useEffect } from 'react';
import { useResource } from '../../custom/api';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils';

import MyPlugin from './MyPlugin';

export interface IMassaPlugin {
  name: string;
  author: string;
  description: string;
  logo: string;
  version: string;
  status: string;
  updatable: boolean;
}

function MyStation() {
  const navigate = useNavigate();
  const {
    error,
    data: plugins,
    status,
  } = useResource<IMassaPlugin[]>('plugin-manager');
  const isLoadingData = status === 'loading';

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  });

  return (
    <>
      <div className="mas-menu-active mb-4 text-neutral">My Station</div>
      {isLoadingData ? (
        <p className="mas-menu-active mb-4 text-neutral">Loading</p>
      ) : plugins && plugins.length > 0 ? (
        <div className="min-w-full flex gap-2 flex-wrap">
          {plugins.map((plugin) => (
            <MyPlugin key={plugin.name} plugin={plugin} />
          ))}
        </div>
      ) : (
        <div className="mas-menu-active mb-4 text-neutral">
          Browse the store below and manage the plugins you've installed in this
          section
        </div>
      )}
    </>
  );
}

export default MyStation;
