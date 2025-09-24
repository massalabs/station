/* eslint-disable max-len */
import Intl from '@/i18n/i18n';
import { motion } from 'framer-motion';
import { 
  FiGlobe, 
  FiSearch, 
  FiUpload, 
  FiBookOpen, 
  FiExternalLink,
  FiArrowRight
} from 'react-icons/fi';
import DeWebLogo from '@/assets/dashboard/deweb.svg';
import { routeFor } from '@/utils/utils';
import { PAGES } from '@/const/pages/pages';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Spinner, useResolveDeweb } from '@massalabs/react-ui-kit';
import { useNetworkStore } from '@/store/store';

interface DeWebComponentProps {
  title: string;
  description: string;
  button: string;
  url: string;
  icon: React.ReactNode;
  features?: string[];
  topics?: string[];
  isLoading?: boolean;
}

interface GetStartedOption {
  text: string;
  url: string;
  isLoading?: boolean;
}

function DeWebComponent({ title, description, button, url, icon, features, topics, isLoading }: DeWebComponentProps) {
  return (
    <motion.div
      className="group relative bg-gradient-to-br from-primary to-primary/90 border border-c-default/50 rounded-2xl p-8 hover:border-c-primary/60 hover:shadow-2xl hover:shadow-c-primary/10 transition-all duration-300 cursor-pointer overflow-hidden backdrop-blur-sm"
      whileHover={{ scale: 1.02, y: -4 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
      onClick={() => {
        if (!isLoading) {
          window.open(url, '_blank', 'noopener,noreferrer');
        }
      }}
    >
      {isLoading ? (
        <div className="flex items-center justify-center h-40">
          <Spinner />
        </div>
      ) : (
        <>
          {/* Gradient overlay on hover */}
          <div className="absolute inset-0 bg-gradient-to-br from-c-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

          <div className="relative z-10">
            <div className="flex items-start gap-6 mb-6">
              <div className="flex-shrink-0 p-4 bg-gradient-to-br from-c-primary/10 to-c-primary/20 rounded-2xl group-hover:from-c-primary/20 group-hover:to-c-primary/30 transition-all duration-300">
                <div className="text-c-primary text-3xl group-hover:scale-110 transition-transform duration-300">
                  {icon}
                </div>
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="mas-subtitle text-neutral mb-3 group-hover:text-c-primary transition-colors duration-300 cursor-default">{title}</h3>
                <p className="mas-body2 text-neutral/80 leading-relaxed mb-6 cursor-default">{description}</p>

                {features && (
                  <div className="mb-6">
                    <ul className="space-y-3">
                      {features.map((feature, index) => (
                        <motion.li
                          key={index}
                          className="flex items-center gap-3 text-sm text-neutral/90"
                          initial={{ opacity: 0, x: -10 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.1 }}
                        >
                          <div className="w-2 h-2 bg-gradient-to-r from-c-primary to-c-primary/80 rounded-full flex-shrink-0 group-hover:scale-125 transition-transform duration-300"></div>
                          <span className="group-hover:text-neutral transition-colors duration-300 cursor-default">{feature}</span>
                        </motion.li>
                      ))}
                    </ul>
                  </div>
                )}

                {topics && (
                  <div className="mb-6">
                    <p className="text-sm text-neutral/70 mb-3 font-medium cursor-default">Key Topics:</p>
                    <div className="flex flex-wrap gap-2">
                      {topics.map((topic, index) => (
                        <motion.span
                          key={index}
                          className="px-3 py-1.5 bg-gradient-to-r from-secondary to-secondary/80 text-xs text-neutral/90 rounded-full border border-c-default/30 hover:border-c-primary/50 hover:from-c-primary/10 hover:to-c-primary/5 transition-all duration-200"
                          initial={{ opacity: 0, scale: 0.8 }}
                          animate={{ opacity: 1, scale: 1 }}
                          transition={{ delay: index * 0.05 }}
                          whileHover={{ scale: 1.05 }}
                        >
                          {topic}
                        </motion.span>
                      ))}
                    </div>
                  </div>
                )}

                <motion.div
                  className="flex justify-end"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <a
                    href={url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-2 px-6 py-3 rounded-xl font-medium transition-all duration-300
                               bg-transparent text-brand border border-brand
                               hover:bg-brand/10
                               hover:shadow-lg hover:shadow-brand/20
                               focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/50"
                    onClick={(e) => e.stopPropagation()}
                  >
                    {button}
                    <FiExternalLink size={18} className="group-hover:translate-x-1 transition-transform duration-300" />
                  </a>
                </motion.div>
              </div>
            </div>
          </div>

          {/* Subtle border gradient */}
          <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-c-primary/20 via-transparent to-c-primary/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-10 blur-xl"></div>
        </>
      )}
    </motion.div>
  );
}

function BenefitCard({ title, description }: { title: string; description: string }) {
  return (
    <motion.div
      className="group relative bg-gradient-to-br from-primary to-primary/95 border border-c-default/50 rounded-2xl p-6 hover:border-c-primary/60 hover:shadow-xl hover:shadow-c-primary/10 transition-all duration-300 overflow-hidden"
      whileHover={{ scale: 1.03, y: -2 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
    >
      {/* Gradient accent */}
      <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-c-primary/60 to-c-primary/30 rounded-t-2xl" />

      {/* Hover overlay */}
      <div className="absolute inset-0 bg-gradient-to-br from-c-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

      <div className="relative z-10">
        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-c-primary/10 to-c-primary/20 rounded-xl flex items-center justify-center group-hover:from-c-primary/20 group-hover:to-c-primary/30 transition-all duration-300">
            <div className="w-6 h-6 bg-gradient-to-r from-c-primary to-c-primary/80 rounded-full"></div>
          </div>
          <div className="flex-1">
            <h4 className="mas-menu-active text-neutral mb-3 group-hover:text-c-primary transition-colors duration-300 cursor-default">{title}</h4>
            <p className="mas-body2 text-neutral/80 leading-relaxed group-hover:text-neutral/90 transition-colors duration-300 cursor-default">{description}</p>
          </div>
        </div>
      </div>

      {/* Subtle glow effect */}
      <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-c-primary/10 via-transparent to-c-primary/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-10 blur-lg"></div>
    </motion.div>
  );
}

function GetStartedSection({ title, description, options }: { title: string; description: string; options: GetStartedOption[] }) {
  return (
    <motion.div
      className="group relative bg-gradient-to-br from-primary to-primary/95 border border-c-default/50 rounded-3xl p-8 hover:border-c-primary/60 hover:shadow-2xl hover:shadow-c-primary/10 transition-all duration-300 overflow-hidden"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      {/* Gradient accent */}
      <div className="absolute top-0 left-0 w-full h-1.5 bg-gradient-to-r from-c-primary via-c-primary/80 to-c-primary/60 rounded-t-3xl" />

      {/* Background pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0 bg-gradient-to-br from-c-primary to-transparent"></div>
      </div>

      <div className="relative z-10">
        <div className="text-center mb-8">
          <h3 className="mas-subtitle text-neutral mb-3 group-hover:text-c-primary transition-colors duration-300 cursor-default">{title}</h3>
          <p className="mas-body2 text-neutral/80 leading-relaxed group-hover:text-neutral/90 transition-colors duration-300 cursor-default">{description}</p>
        </div>

        <div className="space-y-4">
          {options.map((option, index) => (
            <motion.div
              key={index}
              className="group/item flex items-start gap-4 p-4 bg-secondary/30 rounded-2xl border border-c-default/20 hover:border-c-primary/40 hover:bg-secondary/50 transition-all duration-200 cursor-pointer"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.1 }}
              whileHover={{ scale: 1.02, x: 4 }}
              onClick={() => {
                if (option.isLoading) return;
                if (option.url.startsWith('/')) {
                  // Internal navigation
                  window.location.href = option.url;
                } else {
                  // External navigation
                  window.open(option.url, '_blank', 'noopener,noreferrer');
                }
              }}
            >
              {option.isLoading ? (
                <div className="flex items-center justify-center w-full py-2">
                  <Spinner />
                </div>
              ) : (
                <>
                  <div className="flex-shrink-0 w-8 h-8 bg-gradient-to-br from-c-primary/20 to-c-primary/10 rounded-xl flex items-center justify-center group-hover/item:from-c-primary/30 group-hover/item:to-c-primary/20 transition-all duration-300">
                    <FiArrowRight size={16} className="text-c-primary group-hover/item:translate-x-1 transition-transform duration-300" />
                  </div>
                  <p className="text-sm text-neutral/90 leading-relaxed group-hover/item:text-neutral transition-colors duration-300 flex-1 cursor-default">{option.text}</p>
                  {!option.url.startsWith('/') && (
                    <div className="flex-shrink-0 opacity-0 group-hover/item:opacity-100 transition-opacity duration-300">
                      <FiExternalLink size={14} className="text-c-primary/70" />
                    </div>
                  )}
                </>
              )}
            </motion.div>
          ))}
        </div>
      </div>

      {/* Glow effect */}
      <div className="absolute inset-0 rounded-3xl bg-gradient-to-r from-c-primary/10 via-transparent to-c-primary/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 -z-10 blur-2xl"></div>
    </motion.div>
  );
}

export function Deweb() {
  /**
   * TODO: Enable dynamic URL resolution when resolveDeweb function is fixed
   * 
   * Uncomment the lines below to enable dynamic DeWeb URL resolution:
   * 1. Uncomment the useResolveDeweb import at the top
   * 2. Uncomment all the useResolveDeweb hook calls below
   * 3. Replace static i18n URLs with resolved URLs in the component objects
   * 4. Update the loading state condition to use isLoadingAnyUrl
   */
  
  const navigate = useNavigate();

  const [getChainId] = useNetworkStore((state) => [state.getChainId]);
  const chainid = getChainId();

  const mnsUrl = useResolveDeweb(Intl.t('deweb.components.mns.url'), chainid);
  const explorerUrl = useResolveDeweb(Intl.t('deweb.components.explorer.url'), chainid);
  const uploaderUrl = useResolveDeweb(Intl.t('deweb.components.uploader.url'), chainid);
  const documentationUrl = useResolveDeweb(Intl.t('deweb.components.documentation.url'), chainid);
  const dewebSearch = useResolveDeweb(Intl.t('deweb.get-started.users.urls.search'), chainid);
  const dewebExplore = useResolveDeweb(Intl.t('deweb.get-started.users.urls.explore'), chainid);
  const mnsDevUrl = useResolveDeweb(Intl.t('deweb.get-started.developers.urls.mns'), chainid);
  const dewebCliDocs = useResolveDeweb(Intl.t('deweb.get-started.developers.urls.cli_docs'), chainid);
  const helloDappDocs = useResolveDeweb(Intl.t('deweb.get-started.developers.urls.hello_dapp'), chainid);

  // Navigate to error page if any resolver reports an error
  const mnsConversionError = mnsUrl.error || explorerUrl.error || uploaderUrl.error || 
                      documentationUrl.error || dewebSearch.error || dewebExplore.error || 
                      mnsDevUrl.error || dewebCliDocs.error || helloDappDocs.error;

  useEffect(() => {
    if (mnsConversionError) {
      navigate(routeFor('error'), { 
        replace: true,
        state: {
          errorData: {
            title: Intl.t('deweb.get-started.mnsProviderResolveError'),
            message: mnsConversionError
          }
        }
      });
    }
  }, [mnsConversionError, navigate]);

  // Define the data structure with proper typing
  const dewebComponents: DeWebComponentProps[] = [
    {
      title: Intl.t('deweb.components.mns.title'),
      description: Intl.t('deweb.components.mns.description'),
      button: Intl.t('deweb.components.mns.button'),
      url: mnsUrl.resolvedUrl,
      icon: <FiGlobe size={24} />,
      features: [
        Intl.t('deweb.components.mns.features.feature1'),
        Intl.t('deweb.components.mns.features.feature2'),
        Intl.t('deweb.components.mns.features.feature3')
      ],
      isLoading: mnsUrl.isLoading,
    },
    {
      title: Intl.t('deweb.components.explorer.title'),
      description: Intl.t('deweb.components.explorer.description'),
      button: Intl.t('deweb.components.explorer.button'),
      url: explorerUrl.resolvedUrl,
      icon: <FiSearch size={24} />, 
      isLoading: explorerUrl.isLoading,
    },
    {
      title: Intl.t('deweb.components.uploader.title'),
      description: Intl.t('deweb.components.uploader.description'),
      button: Intl.t('deweb.components.uploader.button'),
      url: uploaderUrl.resolvedUrl,
      icon: <FiUpload size={24} />,
      features: [
        Intl.t('deweb.components.uploader.features.feature1'),
        Intl.t('deweb.components.uploader.features.feature2'),
        Intl.t('deweb.components.uploader.features.feature3'),
        Intl.t('deweb.components.uploader.features.feature4')
      ],
      isLoading: uploaderUrl.isLoading,
    },
    {
      title: Intl.t('deweb.components.documentation.title'),
      description: Intl.t('deweb.components.documentation.description'),
      button: Intl.t('deweb.components.documentation.button'),
      url: documentationUrl.resolvedUrl,
      icon: <FiBookOpen size={24} />,
      topics: [
        Intl.t('deweb.components.documentation.topics.topic1'),
        Intl.t('deweb.components.documentation.topics.topic2'),
        Intl.t('deweb.components.documentation.topics.topic3'),
        Intl.t('deweb.components.documentation.topics.topic4')
      ],
      isLoading: documentationUrl.isLoading,
    }
  ];

  const benefits = [
    {
      title: Intl.t('deweb.benefits.items.item0.title'),
      description: Intl.t('deweb.benefits.items.item0.description')
    },
    {
      title: Intl.t('deweb.benefits.items.item1.title'),
      description: Intl.t('deweb.benefits.items.item1.description')
    },
    {
      title: Intl.t('deweb.benefits.items.item2.title'),
      description: Intl.t('deweb.benefits.items.item2.description')
    },
    {
      title: Intl.t('deweb.benefits.items.item3.title'),
      description: Intl.t('deweb.benefits.items.item3.description')
    }
  ];

  const userOptions: GetStartedOption[] = [
    {
      text: Intl.t('deweb.get-started.users.options.option0'),
      url: routeFor(PAGES.STORE),
      isLoading: false,
    },
    {
      text: Intl.t('deweb.get-started.users.options.option1'),
      url: dewebSearch.resolvedUrl,
      isLoading: dewebSearch.isLoading,
    },
    {
      text: Intl.t('deweb.get-started.users.options.option2'),
      url: dewebExplore.resolvedUrl,
      isLoading: dewebExplore.isLoading,
    }
  ];

  const developerOptions: GetStartedOption[] = [
    {
      text: Intl.t('deweb.get-started.developers.options.option0'),
      url: mnsDevUrl.resolvedUrl,
      isLoading: mnsDevUrl.isLoading,
    },
    {
      text: Intl.t('deweb.get-started.developers.options.option1'),
      url: dewebCliDocs.resolvedUrl,
      isLoading: dewebCliDocs.isLoading,
    },
    {
      text: Intl.t('deweb.get-started.developers.options.option2'),
      url: helloDappDocs.resolvedUrl,
      isLoading: helloDappDocs.isLoading,
    },
    {
      text: Intl.t('deweb.get-started.developers.options.option3'),
      url: helloDappDocs.resolvedUrl,
      isLoading: helloDappDocs.isLoading,
    }
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary via-primary to-primary/95 text-neutral relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 opacity-30">
        <div className="absolute top-20 left-10 w-72 h-72 bg-c-primary/20 rounded-full blur-3xl"></div>
        <div className="absolute bottom-20 right-10 w-96 h-96 bg-c-primary/10 rounded-full blur-3xl"></div>
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-gradient-to-r from-c-primary/5 via-transparent to-c-primary/5 rounded-full blur-3xl"></div>
      </div>

      {/* Hero Section */}
      <div className="relative z-10 text-center py-24 px-4">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, ease: "easeOut" }}
        >
          <motion.div
            className="flex justify-center mb-8"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ duration: 0.6, delay: 0.2 }}
          >
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-brand/20 to-brand/20 rounded-full blur-3xl scale-150"></div>

              {/* DeWeb logo with brand color gradient */}
              <div
                className={`relative p-8 bg-gradient-to-br from-brand/90 to-brand/95 rounded-3xl backdrop-blur-sm border border-brand/30 items-center justify-center flex`}
                style={{
                  WebkitMaskImage: `url(${DeWebLogo})`,
                  maskImage: `url(${DeWebLogo})`,
                  WebkitMaskSize: '150px 46px',
                  maskSize: '150px 46px',
                  WebkitMaskPosition: 'center',
                  maskPosition: 'center',
                  WebkitMaskRepeat: 'no-repeat',
                  maskRepeat: 'no-repeat',
                  width: '200px',
                  height: '80px',
                  background: 'linear-gradient(90deg, rgb(0 254 109) 0%, rgb(0 254 109) 100%)',
                  filter: 'drop-shadow(0 0 8px rgba(0,254,109,0.4)) drop-shadow(0 0 16px rgba(0,254,109,0.2))'
                }}
              />
            </div>
          </motion.div>

          <motion.h1
            className="mas-h1 text-neutral mb-6 bg-gradient-to-r from-neutral via-neutral to-neutral/80 bg-clip-text text-transparent cursor-default"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.4 }}
          >
            {Intl.t('deweb.title')}
          </motion.h1>

          <motion.p
            className="mas-h2 text-neutral/90 mb-8 cursor-default"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.6 }}
          >
            {Intl.t('deweb.subtitle')}
          </motion.p>

          <motion.p
            className="mas-body text-neutral/80 max-w-4xl mx-auto leading-relaxed cursor-default"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.8 }}
          >
            {Intl.t('deweb.description')}
          </motion.p>
        </motion.div>
      </div>

      {/* Vision Section */}
      <motion.div
        className="relative z-10 px-4 mb-20"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.2 }}
      >
        <div className="max-w-5xl mx-auto">
          <motion.div
            className="relative bg-gradient-to-br from-primary/80 to-primary/90 backdrop-blur-sm border border-c-primary/30 rounded-3xl p-12 shadow-2xl shadow-c-primary/10"
            whileHover={{ scale: 1.01 }}
            transition={{ duration: 0.3 }}
          >
            {/* Decorative elements */}
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-c-primary via-c-primary/80 to-c-primary rounded-t-3xl"></div>
            <div className="absolute -top-4 left-1/2 transform -translate-x-1/2 w-8 h-8 bg-gradient-to-r from-c-primary to-c-primary/80 rounded-full border-4 border-primary"></div>

            <div className="text-center">
              <motion.h2
                className="mas-h2 text-neutral mb-6 bg-gradient-to-r from-neutral via-c-primary/80 to-neutral bg-clip-text text-transparent cursor-default"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.3 }}
              >
                {Intl.t('deweb.vision.title')}
              </motion.h2>
              <motion.p
                className="mas-body text-neutral/80 text-center leading-relaxed text-lg cursor-default"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.5 }}
              >
                {Intl.t('deweb.vision.description')}
              </motion.p>
            </div>

            {/* Subtle glow */}
            <div className="absolute inset-0 rounded-3xl bg-gradient-to-r from-c-primary/10 via-transparent to-c-primary/10 opacity-50 blur-2xl -z-10"></div>
          </motion.div>
        </div>
      </motion.div>

      {/* Ecosystem Components */}
      <motion.div
        className="relative z-10 px-4 mb-24"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.4 }}
      >
        <div className="max-w-7xl mx-auto">
          <motion.div
            className="text-center mb-16"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.5 }}
          >
            <h2 className="mas-h2 text-neutral mb-6 bg-gradient-to-r from-neutral via-c-primary/80 to-neutral bg-clip-text text-transparent cursor-default">
              {Intl.t('deweb.components.title')}
            </h2>
            <p className="mas-body text-neutral/80 max-w-3xl mx-auto leading-relaxed mb-6 cursor-default">
              {Intl.t('deweb.components.description')}
            </p>
          </motion.div>

          <div className="grid lg:grid-cols-2 gap-8">
            {dewebComponents.map((component, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
              >
                <DeWebComponent {...component} />
              </motion.div>
            ))}
          </div>
        </div>
      </motion.div>

      {/* Benefits Section */}
      <motion.div
        className="relative z-10 px-4 mb-24"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.6 }}
      >
        <div className="max-w-7xl mx-auto">
          <motion.div
            className="text-center mb-16"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.7 }}
          >
            <h2 className="mas-h2 text-neutral mb-6 bg-gradient-to-r from-neutral via-c-primary/80 to-neutral bg-clip-text text-transparent cursor-default">
              {Intl.t('deweb.benefits.title')}
            </h2>
            <div className="w-24 h-1 bg-gradient-to-r from-c-primary to-c-primary/60 mx-auto rounded-full"></div>
          </motion.div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {benefits.map((benefit, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: index * 0.1 + 0.8 }}
              >
                <BenefitCard {...benefit} />
              </motion.div>
            ))}
          </div>
        </div>
      </motion.div>

      {/* Get Started Section */}
      <motion.div
        className="relative z-10 px-4 mb-24"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.8 }}
      >
        <div className="max-w-7xl mx-auto">
          <motion.div
            className="text-center mb-16"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.9 }}
          >
            <h2 className="mas-h2 text-neutral mb-6 bg-gradient-to-r from-neutral via-c-primary/80 to-neutral bg-clip-text text-transparent cursor-default">
              {Intl.t('deweb.get-started.title')}
            </h2>
            <p className="mas-body text-neutral/80 max-w-3xl mx-auto leading-relaxed mb-6 cursor-default">
              {Intl.t('deweb.get-started.description')}
            </p>
          </motion.div>

          <div className="grid lg:grid-cols-2 gap-10">
            <motion.div
              initial={{ opacity: 0, x: -30 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 1.0 }}
            >
              <GetStartedSection
                title={Intl.t('deweb.get-started.users.title')}
                description={Intl.t('deweb.get-started.users.description')}
                options={userOptions}
              />
            </motion.div>
            <motion.div
              initial={{ opacity: 0, x: 30 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 1.1 }}
            >
              <GetStartedSection
                title={Intl.t('deweb.get-started.developers.title')}
                description={Intl.t('deweb.get-started.developers.description')}
                options={developerOptions}
              />
            </motion.div>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
