import explorer from '../../../assets/dashboard/Explorer.svg';
import { RedirectTile } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';

export function Explorer() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://explorer.massa.net/"
        className="bg-brand hover:bg-brand/90 rounded-lg border border-brand h-full relative overflow-hidden"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-brand/90 to-brand"></div>
        <div className="relative z-10 h-full p-8">
          <div className="flex flex-col h-full justify-between gap-6">
            <motion.img
              src={explorer}
              alt="Massa Explorer"
              className="w-14 h-14"
              whileHover={{ 
                scale: 1.1,
                rotate: 5,
                transition: { duration: 0.3 }
              }}
            />
            <div className="flex items-end justify-between">
              <div className="text-c-default font-bold text-4xl leading-tight cursor-default">
                <div>Massa</div>
                <div>Explorer</div>
              </div>
            </div>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
