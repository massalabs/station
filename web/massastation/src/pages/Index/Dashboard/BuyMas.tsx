import { RedirectTile } from '@massalabs/react-ui-kit';
import { motion } from 'framer-motion';
import Intl from '@/i18n/i18n';

export function BuyMas() {
  return (
    <motion.div className="h-full" whileHover={{ scale: 1.03 }}>
      <RedirectTile
        url="https://www.massa.net/get-mas"
        className="rounded-lg border border-c-dark h-full relative overflow-hidden bg-c-dark"
      >
        <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent"></div>
        <div className="relative z-10 h-full p-8">
          <div className="flex flex-col h-full justify-between items-start">
            <motion.div
              className="w-14 h-14 rounded-full bg-brand/20 border-2 border-brand 
                flex items-center justify-center cursor-default"
              whileHover={{ 
                scale: 1.1,
                rotate: 5,
                transition: { duration: 0.3 }
              }}
            >
              <span className="text-brand font-bold text-xl">$</span>
            </motion.div>
            <div className="text-brand font-bold text-4xl leading-tight text-left cursor-default">
              <div>Buy</div>
              <div>MAS</div>
            </div>
            <p className="text-f-primary font-medium text-xs text-left cursor-default">
              {Intl.t('modules.buy-mas.baseline')}
            </p>
          </div>
        </div>
      </RedirectTile>
    </motion.div>
  );
}
