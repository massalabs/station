import { SidePanel } from '@massalabs/react-ui-kit';
import { FiGlobe } from 'react-icons/fi';
import Intl from '../../../i18n/i18n';

export function UploadPlaceholder() {
  return (
    <SidePanel customClass="border-l border-c-default bg-secondary flex ">
      <div className="w-fit p-4">
        <div className="flex flex-col h-full items-center w-[340px] justify-center gap-10 text-neutral">
          <FiGlobe size={114} />
          <p className="mas-title text-center">
            {Intl.t('search.sidebar.placeholder-title')}
          </p>
          <p className="mas-button-active text-center">
            {Intl.t('search.sidebar.placeholder-description')}
          </p>
        </div>
      </div>
    </SidePanel>
  );
}
