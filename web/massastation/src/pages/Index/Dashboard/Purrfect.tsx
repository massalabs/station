import purrfectCat from '../../../assets/dashboard/PurrfectCat.svg';
import Intl from '@/i18n/i18n';
import { RedirectTile } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';

export function Purrfect() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.05 }}>
      <RedirectTile
        size="cs"
        customSize="h-full"
        customClass="bg-neutral text-neutral hover:bg-[#DADADA] hover:cursor-pointer"
        url="https://www.purrfectuniverse.com/"
        className="bg-neutral text-primary rounded-md p-0 hover:bg-c-hover hover:cursor-pointer h-full"
      >
        <div
          style={{
            backgroundImage: `url(${purrfectCat})`,
            backgroundSize: 'cover',
            backgroundRepeat: 'no-repeat',
          }}
          className="h-full flex flex-col gap-4 p-4"
        >
          <p className="mas-subtitle">{Intl.t('modules.purrfect')}</p>
          <p>{Intl.t('dashboard.purrfect-desc')}</p>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
