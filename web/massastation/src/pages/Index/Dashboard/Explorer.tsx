import explorer from '../../../assets/dashboard/Explorer.svg';
import { RedirectTile } from '@massalabs/react-ui-kit';
import { FiSearch } from 'react-icons/fi';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';

export function Explorer() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.05 }}>
      <RedirectTile
        size="cs"
        url="https://explorer.massa.net/"
        className="bg-brand rounded-md p-0 hover:bg-c-hover hover:cursor-pointer h-full"
      >
        <div className="flex flex-col justify-end h-full">
          <div className="flex items-center gap-2 p-4 ml-4 h-8 bg-secondary w-48 rounded-t-md">
            <img width={20} height={20} src={explorer} alt="Massa Explorer" />
            <p className="mas-body2">Massa Explorer</p>
          </div>
          <div className="bg-secondary h-[50%] w-full rounded-b-md p-4">
            <div className="flex items-center justify-around p-2 bg-primary h-full rounded-md">
              <p className="mas-caption">{Intl.t('dashboard.explorer-desc')}</p>
              <FiSearch size={16} />
            </div>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
