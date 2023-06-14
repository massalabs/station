import { FiCodepen, FiGlobe, FiHome } from 'react-icons/fi';
import { Navigator, LayoutStation } from '@massalabs/react-ui-kit';
import MyStation from './MyStation';

export function Store() {
  let navigator = (
    <Navigator
      items={[
        {
          icon: <FiHome />,
          isActive: false,
        },
        {
          icon: <FiCodepen />,
          isActive: true,
        },
        {
          icon: <FiGlobe />,
          isActive: false,
        },
      ]}
      onClickNext={() => console.log('Next clicked')}
      onClickBack={() => console.log('Back clicked')}
    />
  );
  return (
    <LayoutStation
      navigator={navigator}
      onSetTheme={(theme: string) => {
        console.log('selected theme', theme);
      }}
    >
      <div className="mt-24 mas-banner text-neutral">Modules</div>
      <div className="mb-6 mt-10">
        <MyStation />
      </div>
    </LayoutStation>
  );
}
