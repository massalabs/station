import { Button, Input, SidePanel } from '@massalabs/react-ui-kit';
import Intl from '../../i18n/i18n';
import { useEffect, useState } from 'react';
import { isZipFile } from '../../utils/massaConstants';
import { usePost } from '../../custom/api';
import { UseQueryResult } from '@tanstack/react-query';
import { IMassaPlugin } from './StationSection/StationSection';

function Install({
  getPlugins,
}: {
  getPlugins: UseQueryResult<IMassaPlugin[]>;
}) {
  const { refetch: refetchPlugins } = getPlugins;
  const [url, setUrl] = useState<string>('');
  const [error, setError] = useState<string>('');
  const {
    mutate,
    error: installError,
    isLoading,
    isSuccess,
  } = usePost('plugin-manager');

  function validate(url: string) {
    setError('');

    if (!isZipFile.test(url)) {
      setError(Intl.t('store.sidepanel.fe-error'));
      return false;
    }

    return true;
  }

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!validate(url)) return;
    const params = { source: url };
    mutate({ params });
  }

  function handleChange(e: any) {
    setError('');
    setUrl(e.target.value);
  }

  useEffect(() => {
    if (installError) {
      setError(Intl.t('store.sidepanel.be-error'));
    }
  }, [installError]);

  useEffect(() => {
    if (isSuccess) {
      refetchPlugins();
    }
  }, [isSuccess]);

  return (
    <SidePanel customClass="border-l border-c-default">
      <form onSubmit={handleSubmit}>
        <div className="flex h-full w-full items-center justify-center">
          <div
            className="flex flex-col justify-center w-[370px] h-fit p-8
                      bg-primary border-dashed border-2 border-c-default
                      "
          >
            <div className="mas-body text-neutral mb-6">
              {Intl.t('store.sidepanel.banner')}
            </div>
            <div className="mas-body2 text-neutral mb-6">
              {Intl.t('store.sidepanel.description', {
                url: 'https://massa_make_plugin.com',
              })}
            </div>
            <div className="bg-secondary p-4">
              <div className="mas-menu-active text-neutral mb-3">
                {Intl.t('store.sidepanel.title')}
              </div>
              <div className="mas-caption text-neutral mb-3">
                {Intl.t('store.sidepanel.subtitle')}
              </div>
              {isLoading ? (
                <div className="mas-body2 mb-4 text-neutral">
                  {Intl.t('store.loading')}
                </div>
              ) : (
                <>
                  <Input
                    placeholder={Intl.t('store.sidepanel.placeholder')}
                    name="url"
                    value={url}
                    onChange={(e) => handleChange(e)}
                    error={error}
                    customClass="bg-primary mb-3"
                  />
                  <Button type="submit" customClass="mt-3">
                    {Intl.t('store.sidepanel.button')}
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>
      </form>
    </SidePanel>
  );
}

export default Install;
