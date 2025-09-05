import { RedirectTile } from '@massalabs/react-ui-kit';
import massaGreen from '../../../assets/dashboard/Massa_green.svg';
import { motion } from 'framer-motion';

export function Massa() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://massa.net/"
        className="rounded-lg border border-c-dark h-full relative overflow-hidden bg-c-dark"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent"></div>
        <div className="relative z-10 h-full p-8 flex flex-col items-center justify-center gap-4">
          <motion.img
            src={massaGreen}
            alt="Massa"
            className="max-w-full max-h-16 object-contain"
            whileHover={{ 
              scale: 1.1,
              transition: { duration: 0.3 }
            }}
          />
          <div className="text-brand font-bold text-2xl leading-tight text-center">
            MASSA
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
