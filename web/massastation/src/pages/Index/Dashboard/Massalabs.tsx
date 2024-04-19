import { RedirectTile } from '@massalabs/react-ui-kit';
import massa from '../../../assets/dashboard/MassaLabs.svg';
import massaNodes from '../../../assets/dashboard/MassaNodes.svg';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';

export function MassaLabs() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        style={{
          backgroundImage: `url(${massaNodes})`,
          backgroundSize: 'cover',
          backgroundRepeat: 'no-repeat',
        }}
        size="cs"
        customSize="h-full"
        url="https://massa.net/"
      >
        <div className="h-fit flex items-center gap-4">
          <div className="flex items-center gap-4 min-w-fit">
            <img src={massa} alt="massa" width={30} height={30} />
            <p className="mas-subtitle">Massa Labs</p>
          </div>
          <p className="mas-body2">{Intl.t('dashboard.massalabs-desc')}</p>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
