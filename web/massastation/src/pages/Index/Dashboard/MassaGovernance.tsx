import { RedirectTile, useResolveDeweb } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';
import { Spinner } from '@massalabs/react-ui-kit';
import { useNetworkStore } from '@/store/store';

export function MassaGovernance() {
  const [getChainId] = useNetworkStore((state) => [state.getChainId]);
  const governanceUrl = useResolveDeweb('https://mip.massa.network/', getChainId());

  if (governanceUrl.isLoading) {
    return (
      <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
        <div className="rounded-lg h-full relative overflow-hidden bg-brand flex items-center justify-center">
          <Spinner />
        </div>
      </motion.div>
    );
  }

  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url={governanceUrl.resolvedUrl}
        className="rounded-lg h-full relative overflow-hidden bg-brand"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-brand/80 to-brand"></div>
        <div className="relative z-10 h-full p-8 flex flex-col justify-between">
          <div className="text-c-default font-bold text-4xl leading-tight text-left cursor-default">
            <div>Massa</div>
            <div>Governance</div>
          </div>
          <p className="text-c-default font-medium text-sm text-left cursor-default">
            {Intl.t('modules.massa-governance.baseline')}
          </p>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
