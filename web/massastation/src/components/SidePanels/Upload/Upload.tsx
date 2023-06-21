import { SyntheticEvent, useEffect, useRef, useState } from 'react';
import Intl from '../../../i18n/i18n';

import {
  Button,
  Input,
  SidePanel,
  TextArea,
  DragDrop,
  Dropdown,
  Identicon,
} from '@massalabs/react-ui-kit';
// import { usePut } from '../../../custom/api';
import {
  validateDescriptionLength,
  validateFileContent,
  validateFileExtension,
  validateWebsiteDescription,
  validateWebsiteName,
} from '../../../validation/upload';
import axios from 'axios';
import { parseForm } from '../../../utils/ParseForm';
import { useResource } from '../../../custom/api';
import { AccountObject } from '../../../models/AccountModel';
import { Loading } from './Loading';

interface IFormError {
  websiteName?: string;
  description?: string;
}

interface IFormObject {
  websiteName: string;
  description: string;
}

// interface IUploadRequest {
// }

// interface IUploadResponse {
// }

export default function Upload() {
  const form = useRef(null);
  const [formError, setFormError] = useState<IFormError | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);
  const [accountsError, setAccountError] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [nickname, setNickname] = useState<string | null>(null);
  // const { mutate: mutableUpload } = usePut<IUploadRequest, IUploadResponse>('websiteUploader/prepare');

  const {
    data: accounts = [],
    error,
    isLoading,
  } = useResource<AccountObject[]>('plugin/massalabs/wallet/api/accounts'); // TODO: declare constants

  useEffect(() => {
    if (error) {
      setAccountError(Intl.t('search.errors.no-accounts'));
    } else {
      setAccountError(null);
    }
  }, [error]);

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

  async function validate(e: SyntheticEvent): Promise<boolean> {
    const formObject = parseForm<IFormObject>(e);
    const { websiteName, description } = formObject;

    // reset errors
    setFileError(null);
    setFormError(null);

    if (nickname === null) {
      setAccountError(Intl.t('search.errors.no-nickname'));
      return false;
    }

    if (websiteName === '') {
      setFormError({ websiteName: Intl.t('search.errors.no-website-name') });
      return false;
    }

    if (!validateWebsiteName(websiteName)) {
      setFormError({
        websiteName: Intl.t('search.errors.invalid-website-name'),
      });
      return false;
    }

    if (description === '') {
      setFormError({ description: Intl.t('search.errors.no-description') });
      return false;
    }

    if (!validateDescriptionLength(description)) {
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

    return true;
  }

  function getDefaultWallet(): string {
    return 'buildnet';
  }

  function uploadWebsite(e: SyntheticEvent) {
    const formObject = parseForm<IFormObject>(e);
    const { websiteName, description } = formObject;
    const bodyFormData = new FormData();
    bodyFormData.append('url', websiteName);
    bodyFormData.append('description', description); // Add the website description to the form data
    bodyFormData.append('nickname', getDefaultWallet());
    bodyFormData.append('zipfile', file as File); // we force compiler type of `file` because we know it's not null
    axios({
      url: `/websiteUploader/prepare`,
      method: 'put',
      data: bodyFormData,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }

  function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();

    validate(e).then((isValid) => {
      if (isValid) {
        uploadWebsite(e);
      }
    });
  }

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
          <div className="bg-secondary rounded-lg p-4 mb-6">
            <p className="mas-menu-active mb-3">
              {Intl.t('search.sidebar.your-account')}
            </p>
            <p className="mas-caption mb-3">
              {Intl.t('search.sidebar.your-account-description')}
            </p>
            <div className="w-64">
              {isLoading ? (
                <Loading />
              ) : (
                <Dropdown options={accountsItems} select={selectedAccountKey} />
              )}
              {accountsError && (
                <p className="mas-body pt-4 text-s-error">{accountsError}</p>
              )}
            </div>
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
            <Button type="submit">{Intl.t('search.buttons.upload')}</Button>
          </form>
          {fileError && (
            <p className="mas-body pt-4 text-s-error">{fileError}</p>
          )}
        </div>
      </div>
    </SidePanel>
  );
}
