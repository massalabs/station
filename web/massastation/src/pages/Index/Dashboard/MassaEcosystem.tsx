import { RedirectTile } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';

export function MassaEcosystem() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://www.massa.net/ecosystem"
        className="rounded-lg border border-c-dark h-full relative overflow-hidden bg-c-dark"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent"></div>
        <div className="relative z-10 h-full p-8">
          <div className="flex flex-col h-full justify-end items-end">
            <div className="text-brand font-bold text-4xl leading-tight text-right">
              <div>Massa</div>
              <div>Ecosystem</div>
            </div>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
