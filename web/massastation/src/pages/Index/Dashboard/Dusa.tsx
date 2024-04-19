import { RedirectTile } from '@massalabs/react-ui-kit';
import dusaWave from '../../../assets/dashboard/DusaWave.svg';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';

export function Dusa() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.05 }}>
      <RedirectTile
        size="cs"
        customSize="h-full"
        url="https://dusa.io/"
        style={{
          backgroundImage: `url(${dusaWave})`,
          backgroundSize: 'cover',
          backgroundRepeat: 'no-repeat',
        }}
      >
        <div className="h-full flex flex-col justify-end gap-4">
          <p className="mas-subtitle">Dusa</p>
          <p>{Intl.t('dashboard.dusa-desc')}</p>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
