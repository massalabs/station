import Intl from '@/i18n/i18n';

import Placeholder from '@/layouts/Placeholder/Placeholder';
import { FiGlobe } from 'react-icons/fi';

export function SearchComingSoon() {
  return (
    <div className="flex justify-center pt-24 text-f-primary">
      <Placeholder
        message={Intl.t('placeholder.teaser-web-on-chain')}
        icon={<FiGlobe size={114} />}
      />
    </div>
  );
}
