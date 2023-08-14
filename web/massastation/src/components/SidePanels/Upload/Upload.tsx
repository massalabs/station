import { SyntheticEvent, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AxiosError } from 'axios';
import Intl from '../../../i18n/i18n';
import { usePut, useResource } from '../../../custom/api';
import { routeFor } from '../../../utils/utils';
import { parseForm } from '../../../utils/ParseForm';
import { useAccountStore } from '../../../store/store';
import { AccountObject } from '../../../models/AccountModel';

import { URL } from '../../../const/url/url';

import {
  Button,
  Input,
  SidePanel,
  TextArea,
  DragDrop,
  Spinner,
  Identicon,
  Dropdown,
} from '@massalabs/react-ui-kit';
import {
  validateDescriptionLength,
  validateFileContent,
  validateFileSize,
  validateFileExtension,
  validateWebsiteDescription,
  validateWebsiteName,
} from '../../../validation/upload';
import { MASSA_WALLET } from '../../../const/const';
import { MassaPluginModel } from '@/models';

interface IFormError {
  websiteName?: string;
  description?: string;
}

interface IFormObject {
  websiteName: string;
  description: string;
}

interface IUploadResponse {
  name: string;
  description: string;
  address: string;
  brokenChunks: string[];
}

interface IUploadError {
  message: string;
}

export default function Upload() {
  const form = useRef(null);
  const navigate = useNavigate();

  const [formError, setFormError] = useState<IFormError | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);
  const [accountsError, setAccountError] = useState<string | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const nickname = useAccountStore((state) => state.currentAccount);
  const setNickname = useAccountStore((state) => state.setCurrentAccount);

  const {
    mutate: mutableUpload,
    isLoading: uploadLoading,
    error: uploadFail,
  } = usePut<FormData, IUploadResponse, AxiosError<IUploadError>>(
    'websiteUploader/prepare',
    {
      'Content-Type': 'multipart/form-data',
    },
  );

  useEffect(() => {
    if (uploadFail) {
      if (
        uploadFail.response &&
        uploadFail.response.data.message.includes(
          'Try another website name, this one is already taken',
        )
      ) {
        setFormError({
          websiteName: Intl.t('search.errors.website-name-already-exists'),
        });
      } else {
        setUploadError(Intl.t('search.errors.upload-error'));
      }
    } else {
      setUploadError(null);
    }
  }, [uploadFail]);

  async function validate(e: SyntheticEvent): Promise<boolean> {
    const formObject = parseForm<IFormObject>(e);
    const { websiteName, description } = formObject;

    setFileError(null);
    setFormError(null);
    setAccountError(null);

    if (!nickname) {
      setAccountError(Intl.t('search.errors.no-nickname'));
      return false;
    }

    if (!websiteName) {
      setFormError({ websiteName: Intl.t('search.errors.no-website-name') });
      return false;
    }

    if (!validateWebsiteName(websiteName)) {
      setFormError({
        websiteName: Intl.t('search.errors.invalid-website-name'),
      });
      return false;
    }

    if (!description) {
      setFormError({ description: Intl.t('search.errors.no-description') });
      return false;
    }

    if (!validateDescriptionLength(description)) {
      // this can never happen because TextArea component has maxLength={280}
      setFormError({
        description: Intl.t('search.errors.description-too-long'),
      });
      return false;
    }

    if (!validateWebsiteDescription(description)) {
      setFormError({
        description: Intl.t('search.errors.invalid-description'),
      });
      return false;
    }

    if (!file) {
      setFileError(Intl.t('search.errors.no-file'));
      return false;
    }

    if (!validateFileExtension(file?.name)) {
      setFileError(Intl.t('search.errors.invalid-file-extension'));
      return false;
    }

    if (!(await validateFileContent(file))) {
      setFileError(Intl.t('search.errors.invalid-file-content'));
      return false;
    }

    if (!validateFileSize(file)) {
      setFileError(Intl.t('search.errors.file-too-big'));
      return false;
    }

    return true;
  }

  function uploadWebsite(e: SyntheticEvent) {
    const { websiteName, description } = parseForm<IFormObject>(e);
    const bodyFormData = new FormData();
    bodyFormData.append('url', websiteName);
    bodyFormData.append('description', description); // Add the website description to the form data
    bodyFormData.append('nickname', nickname as string); // we force compiler type because we know it's not null
    bodyFormData.append('zipfile', file as File); // we force compiler type of `file` because we know it's not null
    mutableUpload(bodyFormData);
  }

  function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();

    validate(e).then((isValid) => {
      if (isValid) {
        uploadWebsite(e);
      }
    });
  }

  const { data: accounts = [] } = useResource<AccountObject[]>(
    `${URL.WALLET_BASE_API}/${URL.WALLET_ACCOUNTS}`,
  );
  const nicknameExistInAccountsList = accounts.find(
    (account) => account.nickname === nickname,
  );

  const accountsItems = accounts.map((account) => ({
    icon: <Identicon username={account.nickname} size={32} />,
    item: account.nickname,
    onClick: () => setNickname(account.nickname),
  }));

  const selectedAccountKey: number = parseInt(
    Object.keys(accounts).find(
      (_, idx) => accounts[idx].nickname === nickname,
    ) || '0',
  );

  const existingAccount: boolean = accounts.length > 0;

  if (!nicknameExistInAccountsList && existingAccount)
    setNickname(accounts[0].nickname);

  const [pluginWalletIsInstalled, setPluginWalletIsInstalled] = useState(false);

  const { data: plugins, isSuccess } =
    useResource<MassaPluginModel[]>('plugin-manager');

  useEffect(() => {
    if (isSuccess) {
      plugins.forEach((plugin) => {
        if (plugin.name === MASSA_WALLET) {
          setPluginWalletIsInstalled(true);
        }
      });
    }
  }, [isSuccess]);

  return (
    <SidePanel customClass="border-l border-c-default bg-secondary">
      <div className="pr-4 m-auto">
        <div
          className={`text-f-primary border-2 border-dashed border-neutral bg-primary p-8`}
        >
          <p className="mas-body mb-6">{Intl.t('search.sidebar.title')}</p>
          <div className="flex gap-3 mb-6">
            <p className="mas-body2">{Intl.t('search.sidebar.how-to')}</p>
            <h3 className="mas-h3 underline cursor-pointer">
              <a href="https://howtouploadwebsite.com" target="_blank">
                howtouploadwebsite.com
              </a>
            </h3>
          </div>
          <div className="mb-6">
            {pluginWalletIsInstalled ? (
              existingAccount ? (
                <Dropdown options={accountsItems} select={selectedAccountKey} />
              ) : (
                <Button
                  onClick={() =>
                    window.open(
                      '/plugin/massa-labs/massa-wallet/web-app/',
                      '_blank',
                    )
                  }
                >
                  {Intl.t('search.buttons.create-account')}
                </Button>
              )
            ) : (
              <Button onClick={() => navigate(routeFor('index'))}>
                {Intl.t('search.buttons.install-wallet')}
              </Button>
            )}
          </div>
          <form ref={form} onSubmit={handleSubmit}>
            <div className="bg-secondary rounded-lg p-4 mb-6">
              <p className="mas-menu-active mb-3">
                {Intl.t('search.sidebar.your-website')}
              </p>
              <p className="mas-caption mb-3">
                {Intl.t('search.sidebar.your-website-description')}
              </p>
              <div className="pb-3">
                <Input
                  defaultValue=""
                  name="websiteName"
                  placeholder={Intl.t('search.inputs.website-name')}
                  customClass="mb-3 bg-primary"
                  error={formError?.websiteName}
                />
              </div>
              <TextArea
                defaultValue=""
                name="description"
                placeholder={Intl.t('search.inputs.website-description')}
                error={formError?.description}
              />
            </div>
            <div className="mb-6">
              <DragDrop
                onFileLoaded={(file) => setFile(file)}
                placeholder={Intl.t('search.inputs.file')}
                allowed={['zip']}
              />
            </div>
            <Button type="submit" disabled={uploadLoading}>
              {uploadLoading && <Spinner variant="button" />}
              {Intl.t('search.buttons.upload')}
            </Button>
          </form>
          {fileError && (
            <p className="mas-body pt-4 text-s-error">{fileError}</p>
          )}
          {uploadLoading && (
            <p className="mas-body pt-4 text-s-info">
              {Intl.t('search.loading')}
            </p>
          )}
          {accountsError && (
            <p className="mas-body pt-4 text-s-error">{accountsError}</p>
          )}
          {uploadError && (
            <p className="mas-body pt-4 text-s-error">{uploadError}</p>
          )}
        </div>
      </div>
    </SidePanel>
  );
}
