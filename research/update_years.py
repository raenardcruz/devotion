import os
import json

def update_years():
    index_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '../src/components/resources/index.json'))
    
    if not os.path.exists(index_path):
        print(f"Error: index.json not found at {index_path}")
        return

    with open(index_path, 'r', encoding='utf-8') as f:
        entries = json.load(f)
        
    years_map = {
        "alexander_of_alexandria": "d. 326 AD",
        "alexander_of_lycopolis": "c. 300 AD",
        "ambrose": "c. 339–397 AD",
        "aphrahat_aphraates": "c. 280–345 AD",
        "apocrypha": "1st–4th c. AD",
        "apology_of_justin_martyr": "c. 155 AD",
        "archelaus": "c. 278 AD",
        "aristides_the_philosopher": "c. 125 AD",
        "arnobius": "c. 255–330 AD",
        "athanasius": "c. 296–373 AD",
        "athenagoras": "c. 133–190 AD",
        "augustine_of_hippo": "354–430 AD",
        "bardesanes": "154–222 AD",
        "barnabas": "c. 70–132 AD",
        "basil_the_great": "329–379 AD",
        "caius": "c. 180–217 AD",
        "clement_of_alexandria": "c. 150–215 AD",
        "clement_of_rome": "c. 35–99 AD",
        "commodianus": "c. 250 AD",
        "councils": "325–787 AD",
        "cyprian_of_carthage": "c. 200–258 AD",
        "cyril_of_jerusalem": "c. 313–386 AD",
        "dionysius_of_rome": "d. 268 AD",
        "dionysius_the_great": "d. 264 AD",
        "ephraim_the_syrian": "c. 306–373 AD",
        "eusebius_of_caesarea": "c. 260–339 AD",
        "gennadius_of_marseilles": "d. c. 496 AD",
        "gregory_nazianzen": "c. 329–390 AD",
        "gregory_of_nyssa": "c. 335–395 AD",
        "gregory_thaumaturgus": "c. 213–270 AD",
        "gregory_the_great_pope": "c. 540–604 AD",
        "hermas": "2nd c. AD",
        "hilary_of_poitiers": "c. 310–367 AD",
        "hippolytus": "c. 170–235 AD",
        "ignatius_of_antioch": "c. 35–108 AD",
        "irenaeus_of_lyons": "c. 130–202 AD",
        "jerome": "c. 347–420 AD",
        "john_cassian": "c. 360–435 AD",
        "john_chrysostom": "c. 347–407 AD",
        "john_of_damascus": "c. 675–749 AD",
        "julius_africanus": "c. 160–240 AD",
        "justin_martyr": "c. 100–165 AD",
        "lactantius": "c. 250–325 AD",
        "leo_the_great_pope": "c. 400–461 AD",
        "letter_of_ignatius": "c. 110 AD",
        "liturgies": "1st–5th c. AD",
        "malchion": "c. 272 AD",
        "mar_jacob": "c. 451–521 AD",
        "mathetes": "c. 130 AD",
        "methodius": "d. c. 311 AD",
        "minucius_felix": "2nd–3rd c. AD",
        "miscellaneous": "various dates",
        "moses_of_chorene": "c. 410–490 AD",
        "novatian": "c. 200–258 AD",
        "origen": "c. 184–253 AD",
        "pamphilus": "c. 240–309 AD",
        "papias": "c. 60–163 AD",
        "peter_of_alexandria": "d. 311 AD",
        "polycarp": "c. 69–155 AD",
        "rufinus": "c. 345–411 AD",
        "socrates_scholasticus": "c. 380–439 AD",
        "sozomen": "c. 400–450 AD",
        "sulpitius_severus": "c. 363–425 AD",
        "tatian": "c. 120–180 AD",
        "tertullian": "c. 155–220 AD",
        "theodoret": "c. 393–458 AD",
        "theodotus": "c. 190 AD",
        "theophilus": "c. 120–190 AD",
        "venantius": "c. 530–609 AD",
        "victorinus": "d. 303 AD",
        "vincent_of_l_rins": "d. c. 445 AD"
    }
    
    updated_count = 0
    for entry in entries:
        doc_id = entry.get('id')
        if doc_id in years_map:
            entry['years'] = years_map[doc_id]
            updated_count += 1
            
    with open(index_path, 'w', encoding='utf-8') as f:
        json.dump(entries, f, indent=2, ensure_ascii=False)
        
    print(f"Successfully updated {updated_count} entries with dates in {index_path}")

if __name__ == '__main__':
    update_years()
