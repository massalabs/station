import { RedirectTile, useResolveDeweb } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';
import syntra from '@/assets/dashboard/Syntra.svg';
import Intl from '@/i18n/i18n';
import { Spinner } from '@massalabs/react-ui-kit';
import { useNetworkStore } from '@/store/store';

export function Syntra() {
  const [getChainId] = useNetworkStore((state) => [state.getChainId]);
  const syntraUrl = useResolveDeweb('https://syntra.massa.network/', getChainId());
  
  if (syntraUrl.isLoading) {
    return (
      <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
        <div className="rounded-lg h-full relative overflow-hidden bg-c-default flex items-center justify-center">
          <Spinner />
        </div>
      </motion.div>
    );
  }

  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url={syntraUrl.resolvedUrl}
        className="rounded-lg h-full relative overflow-hidden bg-c-default"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-c-default/80 to-c-default"></div>
        <div className="relative z-10 h-full p-8 flex flex-col items-center justify-center">
          <motion.img
            src={syntra}
            alt="Syntra"
            className="max-w-full max-h-3/4 object-contain mb-4"
            whileHover={{ scale: 1.05 }}
            transition={{ duration: 0.2 }}
          />
          <p className="text-center text-white font-medium text-sm cursor-default">
            {Intl.t('modules.syntra.baseline')}
          </p>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
