import { RedirectTile } from '@massalabs/react-ui-kit';
import foundation from '../../../assets/dashboard/Foundation.svg';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';
export function Foundation() {
  return (
    <motion.div whileHover={{ scale: 1.05 }}>
      <RedirectTile
        url="https://massa.foundation/"
        className="bg-neutral text-primary rounded-md p-0 hover:bg-c-hover hover:cursor-pointer h-full"
      >
        <div className="h-full flex flex-col gap-4 p-4">
          <p className="mas-subtitle">{Intl.t('modules.foundation')}</p>
          <p>{Intl.t('dashboard.foundation-desc')}</p>
        </div>
        <div className="relative flex justify-end">
          <img
            width={140}
            height={140}
            src={foundation}
            alt={Intl.t('modules.foundation')}
          />
        </div>
      </RedirectTile>
    </motion.div>
  );
}
