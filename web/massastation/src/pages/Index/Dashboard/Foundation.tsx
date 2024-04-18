import { RedirectTile } from '@massalabs/react-ui-kit';
import foundation from '../../../assets/dashboard/Foundation.svg';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';
export function Foundation() {
  return (
    <motion.div whileHover={{ scale: 1.05 }}>
      <RedirectTile
        url="https://massa.foundation/"
        className="bg-neutral text-primary rounded-md p-0 hover:bg-c-hover hover:cursor-pointer h-fit"
      >
        <div className="h-full flex flex-col gap-4 p-4">
          <p className="mas-subtitle">Massa Foundation</p>
          <p>{Intl.t('dashboard.foundation-desc')}</p>
        </div>
        <div className="relative flex justify-end">
          <img
            width={120}
            height={120}
            src={foundation}
            alt="Massa Foundation"
          />
        </div>
      </RedirectTile>
    </motion.div>
  );
}
