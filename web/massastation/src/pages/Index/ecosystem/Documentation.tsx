import { Link } from 'react-router-dom';
import Intl from '@/i18n/i18n';
import { Button, MassaLogo } from '@massalabs/react-ui-kit';
import { FiExternalLink } from 'react-icons/fi';
import { docUrl } from '@/const';

export function Documentation() {
  return (
    <Link target="_blank" to={docUrl}>
      <div
        className="flex flex-col items-center justify-center gap-8 bg-tertiary 
            h-full w-full rounded-2xl border-2 border-secondary p-4"
      >
        <MassaLogo size={52} />
        <Button posIcon={<FiExternalLink />}>{Intl.t('index.docs')}</Button>
      </div>
    </Link>
  );
}
