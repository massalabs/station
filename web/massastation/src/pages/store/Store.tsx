import { MyStation } from './MyStation';

import Intl from '../../i18n/i18n';

export function Store() {
  return (
    <>
      <div className="mas-banner text-neutral mb-10 mt-24">
        {Intl.t('store.modules-banner')}
      </div>
      <div className="mas-menu-active mb-4 text-neutral">
        {Intl.t('store.mystation-banner')}
      </div>
      <MyStation />
    </>
  );
}
