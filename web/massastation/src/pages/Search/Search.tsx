import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { routeFor } from '../../utils/utils';
import Intl from '@/i18n/i18n';

import { useResource } from '@/custom/api';
import { DomainModel } from '@/models/DomainModel';
import Upload from '@/components/SidePanels/Upload/Upload';
import { FAVORITES_WEBSITES } from '@/const';

import { Loading } from './Loading';
import DomainModelItem from './DomainModelItem/DomainModelItem';

export function Search() {
  const navigate = useNavigate();

  const {
    data: websites = [],
    error,
    isSuccess,
    isLoading,
  } = useResource<DomainModel[]>('all/domains');

  const fav: DomainModel[] = [];

  if (isSuccess) {
    websites.forEach((website) => {
      if (FAVORITES_WEBSITES.includes(website.name)) {
        fav.push(website);
      }
    });
  }

  useEffect(() => {
    if (error) {
      navigate(routeFor('error'));
    }
  }, [error, navigate]);

  return (
    <>
      {isLoading ? (
        <Loading />
      ) : (
        <div className="bg-primary text-f-primary pt-24">
          <h1 className="mas-banner mb-10">{Intl.t('search.title-banner')}</h1>
          <div className="overflow-auto h-[65vh]">
            {fav.length > 0 && (
              <>
                <p className="mas-body mb-3">{Intl.t('search.fav-websites')}</p>
                <div className="flex flex-col gap-5 mb-14">
                  {fav.map((fav: DomainModel, index: number) => (
                    <div key={index}>
                      <DomainModelItem website={fav} />
                    </div>
                  ))}
                </div>
              </>
            )}
            <p className="mas-body pb-3">{Intl.t('search.all-websites')}</p>
            <div className="flex flex-col gap-5">
              {websites.map((website: DomainModel, index: number) => (
                <div key={index}>
                  <DomainModelItem website={website} />
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
      <Upload />
    </>
  );
}
