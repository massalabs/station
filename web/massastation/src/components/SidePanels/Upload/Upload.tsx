import Intl from '../../../i18n/i18n';

import {
  Button,
  Input,
  SidePanel,
  TextArea,
  DragDrop,
} from '@massalabs/react-ui-kit';

export default function Upload() {
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
              howtouploadwebsite.com
            </h3>
          </div>
          <div className="bg-secondary rounded-lg p-4 mb-6">
            <p className="mas-menu-active mb-3">
              {Intl.t('search.sidebar.your-website')}
            </p>
            <p className="mas-caption mb-3">
              {Intl.t('search.sidebar.your-website-desc')}
            </p>
            <Input
              placeholder={Intl.t('search.inputs.website-name')}
              customClass="mb-3 bg-primary"
            />
            <TextArea placeholder={Intl.t('search.inputs.website-desc')} />
          </div>
          <div className="mb-6">
            <DragDrop allowed={['zip']} />
          </div>
          <Button>{Intl.t('search.buttons.upload')}</Button>
        </div>
      </div>
    </SidePanel>
  );
}
