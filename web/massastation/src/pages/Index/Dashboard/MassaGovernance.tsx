import { RedirectTile } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';

export function MassaGovernance() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://mip.massa.network/"
        className="rounded-lg h-full relative overflow-hidden bg-brand"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-brand/80 to-brand"></div>
        <div className="relative z-10 h-full p-8 flex items-end justify-start">
          <div className="text-c-default font-bold text-4xl leading-tight text-left">
            <div>Massa</div>
            <div>Governance</div>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
