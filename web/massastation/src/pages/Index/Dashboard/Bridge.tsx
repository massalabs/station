import { RedirectTile } from '@massalabs/react-ui-kit';
import bridgeOutline from '../../../assets/dashboard/BridgeOutline.svg';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';

export function Bridge() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.05 }}>
      <RedirectTile
        url="https://bridge.massa.net"
        size="cs"
        customSize="h-full"
        style={{
          backgroundImage: `url(${bridgeOutline})`,
          backgroundSize: 'cover',
          backgroundRepeat: 'no-repeat',
        }}
      >
        <div className="relative top-10 pl-2 h-fit flex flex-col gap-4 justify-between">
          <p className="mas-subtitle">{Intl.t('modules.bridge')}</p>
          <p>{Intl.t('dashboard.bridge-desc')}</p>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
