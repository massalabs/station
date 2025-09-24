import { Link, useLocation } from 'react-router-dom';
import Intl from '@/i18n/i18n';
import { routeFor } from '@/utils/utils';

export type ErrorData = {
  message: string;
  title: string;
};

export function Error() {
  const location = useLocation();

  const errorData: ErrorData = location.state?.errorData;

  return (
    <div
      id="error-page"
      className="text-f-primary min-h-screen flex flex-col items-center justify-center text-center gap-4 p-4"
    >
      <h1 className="mas-banner cursor-default">{errorData.title || Intl.t('unexpected-error.title')}</h1>
      <p className="mas-body">{errorData.message || Intl.t('unexpected-error.description')}</p>
      <Link to={routeFor('index')} className="underline">
        {Intl.t('unexpected-error.link')}
      </Link>
    </div>
  );
}
