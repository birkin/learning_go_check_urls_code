import datetime, pprint
import requests


sites = {
    'repo_file': {
        'url': 'https://repository.library.brown.edu/storage/bdr:6758/PDF/',
        'expected': 'BleedBox' },
    'repo_search': {
        'url': 'https://repository.library.brown.edu/studio/search/?q=elliptic',
        'expected': 'The sequence of division polynomials' },
    'bipg_wiki': {
        'url': 'https://wiki.brown.edu/confluence/display/bipg/Brown+Internet+Programming+Group+Home',
        'expected': 'The BIPG idea' },
    'booklocator_app': {
        'url': 'http://library.brown.edu/services/book_locator/?callnumber=GC97+.C46&location=sci&title=Chemistry+and+biochemistry+of+estuaries&status=AVAILABLE&oclc_number=05831908&public=true',
        'expected': 'GC97 .C46 Level 11, Aisle 2A' },
    'callnumber_app': {
        'url': 'https://apps.library.brown.edu/callnumber/v2/?callnumber=PS3576',
        'expected': 'American Literature' },
    'clusters api': {
        'url': 'https://library.brown.edu/clusters_api/data/',
        'expected': 'scili-friedman' },
    'easyborrow_feed': {
        'url': 'http://library.brown.edu/easyborrow/feeds/latest_items/',
        'expected': 'easyBorrow -- recent requests' },
    'freecite': {
        'url': 'http://freecite.library.brown.edu/welcome/',
        'expected': 'About FreeCite' },
    'iip_inscriptions': {
        'url': 'http://library.brown.edu/cds/projects/iip/viewinscr/abur0001/',
        'expected': 'Khirbet Abu Rish' },
    'iip_processor': {
        'url': 'https://apps.library.brown.edu/iip_processor/info/',
        'expected': 'hi' },
    }

results = {}

start = datetime.datetime.now()
for ( location, value_dct ) in sites.items():
    url = sites[location]['url']
    expected = sites[location]['expected']
    mini_start = datetime.datetime.now()
    r = requests.get( url )
    result = 'not_found'
    assert type( r.content ) == bytes, type(r.content)
    # print( 'r.content, ```%s```' % r.content )
    if expected.encode('utf-8') in r.content:
        result = 'found'
    now = datetime.datetime.now()
    elapsed = now - mini_start
    print( 'location `%s` took `%s`' % (location, elapsed) )
    results[location] = result
end = datetime.datetime.now()

time_taken = end - start

print( 'results, ```%s```' % pprint.pformat(results) )
print( 'time_taken, `%s`' % time_taken )
