import { useResource } from '../../custom/api';
import StationSection, { IMassaPlugin } from './StationSection/StationSection';
import StoreSection from './StoreSection/StoreSection';
import Intl from '../../i18n/i18n';
import Install from './Install';

export function Store() {
  const getPlugins = useResource<IMassaPlugin[]>('plugin-manager');
  return (
    <>
      <div className="bg-primary text-f-primary pt-24">
        <h1 className="mas-banner mb-10"> {Intl.t('store.modules-banner')}</h1>
        <div className="overflow-auto h-[65vh]">
          <div className="mas-body mb-3 text-neutral">
            {Intl.t('store.mystation-banner')}
          </div>
          <div className="mb-10">
            <StationSection getPlugins={getPlugins} />
          </div>
          <div className="mas-body mb-3 text-neutral">
            {Intl.t('store.store-banner')}
          </div>
          <StoreSection getPlugins={getPlugins} />
        </div>
      </div>
      <Install />
    </>
  );
}
