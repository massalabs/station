import { RedirectTile } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';
import syntra from '@/assets/dashboard/Syntra.svg';

export function Syntra() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://syntra.massa.network"
        className="rounded-lg h-full relative overflow-hidden bg-c-default"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-c-default/80 to-c-default"></div>
        <div className="relative z-10 h-full p-8 flex items-center justify-center">
          <motion.img
            src={syntra}
            alt="Syntra"
            className="max-w-full max-h-full object-contain"
            whileHover={{ scale: 1.05 }}
            transition={{ duration: 0.2 }}
          />
        </div>
      </RedirectTile>
    </motion.div>
  );
}
