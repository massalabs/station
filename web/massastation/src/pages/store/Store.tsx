import { FiCodepen, FiGlobe, FiHome } from 'react-icons/fi';
import { Navigator, LayoutStation } from '@massalabs/react-ui-kit';
import Intl from '../../i18n/i18n';
import StationSection from './StationSection/StationSection';
import StoreSection from './StoreSection/StoreSection';

export function Store() {
  let navigator = (
    <Navigator
      items={[
        {
          icon: <FiHome />,
          isActive: false,
        },
        {
          icon: <FiCodepen />,
          isActive: true,
        },
        {
          icon: <FiGlobe />,
          isActive: false,
        },
      ]}
      // these are mandatory and cannot be remove
      // correct redirect will be implemented later
      onClickNext={() => console.log('Next clicked')}
      onClickBack={() => console.log('Back clicked')}
    />
  );
  return (
    <LayoutStation navigator={navigator}>
      <div className="mb-10">
        <div className="mas-banner text-neutral mb-10 mt-24">
          {Intl.t('store.modules-banner')}
        </div>
        <div className="mas-menu-active mb-4 text-neutral">
          {Intl.t('store.mystation-banner')}
        </div>
        <StationSection />
      </div>
      <div>
        <div className="mas-menu-active mb-4 text-neutral">
          {Intl.t('store.store-banner')}
        </div>
        <StoreSection />
      </div>
    </LayoutStation>
  );
}
