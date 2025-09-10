import { Button } from '@massalabs/react-ui-kit';
import { FiCodepen, FiGlobe } from 'react-icons/fi';
import { useNavigate } from 'react-router-dom';
import { MassaPluginModel, MassaStoreModel } from '@/models';
import Intl from '@/i18n/i18n';
import { routeFor } from '@/utils/utils';
import { useResource } from '../../custom/api';
import { UseQueryResult } from '@tanstack/react-query';
import { DashboardStation } from './DashboardStation';

export function Index() {
  const plugins = useResource<MassaPluginModel[]>('plugin-manager');
  const store = useResource<MassaStoreModel[]>('plugin-store');
  const { isLoading: isPluginsLoading } = plugins;
  const { isLoading: isStoreLoading } = store;

  return isPluginsLoading || isStoreLoading ? (
    <>Loading</>
  ) : (
    <NestedIndex store={store} plugins={plugins}></NestedIndex>
  );
}

function NestedIndex({
  store,
  plugins,
}: {
  store: UseQueryResult<MassaStoreModel[]>;
  plugins: UseQueryResult<MassaPluginModel[]>;
}) {
  const navigate = useNavigate();
  const { data: massaPlugins } = plugins;
  const availablePlugins = store.data;

  return (
    <>
      <div className="bg-primary text-f-primary pt-24">
        <h1 className="mas-banner mb-10 cursor-default"> {Intl.t('index.title-banner')}</h1>
        <div className="w-fit max-w-[1760px] ">
          <div className="flex gap-8 pb-10">
            <Button
              variant="primary"
              preIcon={<FiGlobe />}
              customClass="w-96 bg-c-default text-brand border border-c-default hover:bg-c-hover rounded-lg"
              onClick={() => navigate(routeFor('deweb'))}
            >
              <div className="flex items-center mas-buttons">
                {Intl.t('index.buttons.deweb')}
              </div>
            </Button>
            <Button
              variant="primary"
              preIcon={<FiCodepen />}
              customClass="w-96 bg-c-default text-brand border border-c-default hover:bg-c-hover rounded-lg"
              onClick={() => navigate(routeFor('store'))}
            >
              <div className="flex items-center mas-buttons">
                {Intl.t('index.buttons.explore')}
              </div>
            </Button>
          </div>
          <DashboardStation
            massaPlugins={massaPlugins}
            availablePlugins={availablePlugins}
          />
        </div>
      </div>
    </>
  );
}
