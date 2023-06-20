import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useResource } from '../../custom/api';
import { DomainsModel } from '../../models/DomainsModel';
import Intl from '../../i18n/i18n';

import Upload from '../../components/SidePanels/Upload/Upload';

import { Description } from '@massalabs/react-ui-kit';
import { FiArrowUpRight } from 'react-icons/fi';

const FAVORITES = ['flappy', 'psychedelic', 'flappynathana1'];

export function Search() {
  const navigate = useNavigate();

  const {
    data: websites = [],
    error,
    isSuccess,
  } = useResource<DomainsModel[]>('all/domains');

  const fav: DomainsModel[] = [];

  if (isSuccess) {
    websites.forEach((website) => {
      if (FAVORITES.includes(website.name)) {
        fav.push(website);
      }
    });
  }

  useEffect(() => {
    if (error) {
      navigate('/error');
    }
  }, [error, navigate]);

  return (
    <>
      <div className="bg-primary text-f-primary pt-24">
        <h1 className="mas-banner mb-10">{Intl.t('search.title-banner')}</h1>
        <div className="overflow-auto h-[700px]">
          {fav.length > 0 && (
            <p className="mas-body mb-3">{Intl.t('search.fav-websites')}</p>
          )}
          <div className="flex flex-col gap-5 mb-14">
            {fav.map((fav: DomainsModel, index: number) => (
              <Description
                key={index}
                variant="secondary"
                preIcon={<FiArrowUpRight />}
                title={fav.name}
                website={fav.name + '.massa'}
                description={fav.description}
              />
            ))}
          </div>
          <p className="mas-body pb-3">{Intl.t('search.all-websites')}</p>
          <div className="flex flex-col gap-5">
            {websites.map((website: DomainsModel, index: number) => (
              <Description
                key={index}
                variant="secondary"
                preIcon={<FiArrowUpRight />}
                title={website.name}
                website={website.name + '.massa'}
                description={website.description}
              />
            ))}
          </div>
        </div>
      </div>
      <Upload />
    </>
  );
}
