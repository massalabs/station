import { useEffect, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useConfigStore } from '@/store/store';

import { Theme, Toast } from '@massalabs/react-ui-kit';
import { useLocalStorage } from '@massalabs/react-ui-kit/src/lib/util/hooks/useLocalStorage';
import { FiCodepen, FiGlobe, FiHome, FiSettings } from 'react-icons/fi';
import { DEFAULT_THEME, LayoutStation } from '@/layouts/LayoutStation/LayoutStation';

import { PAGES } from '@/const/pages/pages';
import { THEME_STORAGE_KEY } from '@/const';

interface INavigatorSteps {
  [key: string]: object;
}

const navigatorSteps: INavigatorSteps = {
  index: {
    previous: null,
    next: PAGES.STORE,
  },
  store: {
    previous: PAGES.INDEX,
    next: PAGES.CONFIG,
  },
  config: {
    previous: PAGES.STORE,
    next: PAGES.DEWEB,
  },
  deweb: {
    previous: PAGES.CONFIG,
    next: null,
  },
};

export function Base() {
  const [theme, setThemeStorage] = useLocalStorage<Theme>(
    THEME_STORAGE_KEY,
    DEFAULT_THEME,
  );

  const { pathname } = useLocation();
  const navigate = useNavigate();

  const currentPage = (() => {
    const lastSegment = pathname.split('/').pop() || 'index';
    // Remove .html extension if present
    const cleanPage = lastSegment.replace(/\.html$/, '') || 'index';
    // Map to valid pages or default to index
    const validPages = ['index', 'store', 'config', 'deweb'];
    return validPages.includes(cleanPage) ? cleanPage : 'index';
  })();
  const [activePage, setActivePage] = useState(currentPage);

  const setThemeStore = useConfigStore((s) => s.setTheme);

  // handle theme-dark to theme-dark-v2 migration
  useEffect(() => {
    if (theme === 'theme-dark') {
      setThemeStorage('theme-dark-v2');
    }
  }, []); // run only once on mount

  useEffect(() => {
    setActivePage(currentPage);
  }, [setActivePage, pathname, currentPage]);

  const STEP = navigatorSteps[currentPage] as INavigatorSteps;

  // Get current step index for stepper
  const steps = [PAGES.INDEX, PAGES.STORE, PAGES.CONFIG, PAGES.DEWEB];
  const currentStepIndex = steps.indexOf(activePage);

  const navigator = (
    <div className="flex items-center gap-4">
      {/* Previous Button */}
      <button
        onClick={handlePrevious}
        disabled={!STEP.previous}
        className={`
          w-10 h-10 rounded-lg flex items-center justify-center transition-all
          ${!STEP.previous 
          ? 'bg-c-disabled-1 border-c-disabled-1 opacity-50 cursor-not-allowed' 
          : 'bg-c-default border-c-default hover:opacity-80'
    }
        `}
      >
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" className="text-white">
          <path 
            d="M15 18L9 12L15 6" 
            stroke="currentColor" 
            strokeWidth="2" 
            strokeLinecap="round" 
            strokeLinejoin="round"
          />
        </svg>
      </button>

      {/* Center Section: Icon Stepper */}
      <div className="flex items-center gap-3">
        {/* Icon Stepper */}
        <div className="flex items-center gap-3">
          {steps.map((step, index) => {
            const isActive = index === currentStepIndex;
            const isCompleted = index < currentStepIndex;
            const iconSize = isActive ? 24 : 16;
            
            const handleStepClick = () => {
              setActivePage(step);
              navigate(step);
            };
            
            return (
              <button
                key={index}
                onClick={handleStepClick}
                className="flex items-center justify-center transition-all duration-300 rounded-lg 
                  hover:bg-opacity-10 hover:bg-white focus:outline-none focus:ring-2 
                  focus:ring-brand focus:ring-opacity-50"
                style={{ 
                  width: isActive ? '32px' : '20px', 
                  height: isActive ? '32px' : '20px' 
                }}
                title={`Go to ${step.charAt(0).toUpperCase() + step.slice(1)} page`}
              >
                {step === PAGES.INDEX && (
                  <FiHome 
                    className={`transition-all duration-300 ${
                      isActive 
                        ? 'text-brand' 
                        : isCompleted 
                        ? 'text-brand' 
                        : 'text-c-disabled-1'
                    }`} 
                    size={iconSize} 
                  />
                )}
                {step === PAGES.STORE && (
                  <FiCodepen 
                    className={`transition-all duration-300 ${
                      isActive 
                        ? 'text-brand' 
                        : isCompleted 
                        ? 'text-brand' 
                        : 'text-c-disabled-1'
                    }`} 
                    size={iconSize} 
                  />
                )}
                {step === PAGES.CONFIG && (
                  <FiSettings 
                    className={`transition-all duration-300 ${
                      isActive 
                        ? 'text-brand' 
                        : isCompleted 
                        ? 'text-brand' 
                        : 'text-c-disabled-1'
                    }`} 
                    size={iconSize} 
                  />
                )}
                {step === PAGES.DEWEB && (
                  <FiGlobe 
                    className={`transition-all duration-300 ${
                      isActive 
                        ? 'text-brand' 
                        : isCompleted 
                        ? 'text-brand' 
                        : 'text-c-disabled-1'
                    }`} 
                    size={iconSize} 
                  />
                )}
              </button>
            );
          })}
        </div>
      </div>

      {/* Next Button */}
      <button
        onClick={handleNext}
        disabled={!STEP.next}
        className={`
          w-10 h-10 rounded-lg flex items-center justify-center transition-all
          ${!STEP.next 
          ? 'bg-c-disabled-1 border-c-disabled-1 opacity-50 cursor-not-allowed'
          : 'bg-c-default border-c-default hover:opacity-80'
    }
        `}
      >
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" className="text-white">
          <path d="M9 18L15 12L9 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
      </button>
    </div>
  );

  // Functions
  function handleNext() {
    let { next } = STEP;

    setActivePage(next.toString());
    navigate(next);
  }

  function handlePrevious() {
    let { previous } = STEP;

    setActivePage(previous.toString());
    navigate(previous);
  }

  function handleSetTheme() {
    let toggledTheme: Theme =
      theme === 'theme-dark-v2' ? 'theme-light' : 'theme-dark-v2';

    setThemeStorage(toggledTheme);
    setThemeStore(toggledTheme);
  }

  // Template
  return (
    <div className={theme}>
      <LayoutStation
        navigator={navigator}
        onSetTheme={handleSetTheme}
        storedTheme={theme}
      >
        <Outlet />
        <Toast />
      </LayoutStation>
    </div>
  );
}
