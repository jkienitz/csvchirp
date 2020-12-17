import csv

# Input fuelds
# Channel Number,Receive Frequency,Transmit Frequency,Offset Frequency,Offset Direction,Operating Mode,Name,Show Name,Tone Mode,CTCSS,DCS,Tx Power,Skip,Step,Clock Shift,Bank 1,Bank 2,Bank 3,Bank 4,Bank 5,Bank 6,Bank 7,Bank 8,Bank 9,Bank 10,Comment,Tx Narrow,Pager Enable,

fieldnames = ['Location', 'Name', 'Frequency', 'Duplex', 'Offset', 'Tone', 'rToneFreq', 'cToneFreq', 'DtcsCode', 'DtcsPolarity', 'Mode', 'TStep', 'Skip', 'Comment', 'URCALL', 'RPT1CALL', 'RPT2CALL']
offsetdirection_values = {"Minus": '-', "Plus": '+', "Simplex": '', "Split": ''}
offset_values = {'5.00 MHz': '5.00', '7.85 MHz': '7.85', '8.00 MHz': '8.00', '9.15 MHz': '9.15', '600 kHz': '0.600000', ' ': ''}
tone_values = {'Tone': 'Tone', 'None': ''}
skip_values = {'Off': '', 'Skip': 'S', 'XXX': 'P'}

counter = 0

with open('XCZFreqListv1_01.csv', newline='') as csvfile:
  reader = csv.DictReader(csvfile)
  with open('chirp.csv', 'w', newline='') as csvoutput:
    writer = csv.DictWriter(csvoutput, fieldnames=fieldnames)
    writer.writeheader()
    for row in reader:

      isSplit = False
      if row['Offset Direction']  ==  'Split':
        isSplit = True
      
      duplex = offsetdirection_values[row['Offset Direction']]
      if row['Tx Power'] == 'Low':
	      duplex = 'off'
      elif isSplit:
        duplex = 'split'

      offset = offset_values[row['Offset Frequency']]
      if isSplit:
        offset = row['Transmit Frequency']


      outrow = {'Location': row['Channel Number'], 
        'Name': row['Name'], 
        'Frequency': row['Receive Frequency'], 
        'Duplex': duplex, 
        'Offset': offset, 
        'Tone': tone_values[row['Tone Mode']], 
        'rToneFreq': row['CTCSS'].split()[0], 
        'cToneFreq': '88.5',  # row[''],  # ???
        'DtcsCode': row['DCS'], 
        'DtcsPolarity': 'NN',  # row[''],   # ???
        'Mode': 'FM',          # row[''], 
        'TStep': '5.00',         # row[''], 
        'Skip': skip_values[row['Skip']], 
        'Comment': row['Comment'], 
        'URCALL': '',       # row[''], 
        'RPT1CALL': '',     # row[''], 
        'RPT2CALL': ''     # row['']
        }
      writer.writerow(outrow)
      counter += 1
print(counter)

      