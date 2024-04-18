import { RedirectTile } from '@massalabs/react-ui-kit';
import massa from '../../../assets/dashboard/MassaLabs.svg';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';

export function MassaLabs() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.05 }}>
      <RedirectTile size="cs" customSize="h-full" url="https://massa.net/">
        <div className="h-fit flex items-center gap-4">
          <img width={60} height={60} src={massa} alt="Massa Website" />
          <div className="flex flex-col gap-2">
            <p className="mas-subtitle">Massa Labs</p>
            <p>{Intl.t('dashboard.massalabs-desc')}</p>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
