import { IMassaPlugin, MyStation } from './MyStation';

import Intl from '../../i18n/i18n';
import { useResource } from '../../custom/api';
import StationSection from './StationSection/StationSection';
import StoreSection from './StoreSection/StoreSection';

export function Store() {
  const getPlugins = useResource<IMassaPlugin[]>('plugin-manager');
  return (
    <>
      <div className="mas-banner text-neutral mb-10 mt-24">
        {Intl.t('store.modules-banner')}
      </div>
      <div className="mas-menu-active mb-4 text-neutral">
        {Intl.t('store.mystation-banner')}
      </div>
      <StationSection getPlugins={getPlugins} />
      <div className="mas-menu-active mb-4 text-neutral">
        {Intl.t('store.store-banner')}
      </div>
      <StoreSection />
    </>
  );
}
