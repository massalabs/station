import { FiCodepen, FiGlobe, FiHome } from 'react-icons/fi';
import { Navigator, LayoutStation } from '@massalabs/react-ui-kit';

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
      <div className="text-f-primary">This is the store page</div>
    </LayoutStation>
  );
}
