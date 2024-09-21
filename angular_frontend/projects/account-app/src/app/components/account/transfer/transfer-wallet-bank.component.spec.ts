import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransferWalletBankComponent } from './transfer-wallet-bank.component';

describe('TransferComponent', () => {
  let component: TransferWalletBankComponent;
  let fixture: ComponentFixture<TransferWalletBankComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TransferWalletBankComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TransferWalletBankComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
